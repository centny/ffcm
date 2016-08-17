using io.vty.cswf.doc;
using io.vty.cswf.log;
using io.vty.cswf.netw.dtm;
using io.vty.cswf.netw.rc;
using io.vty.cswf.netw.sck;
using io.vty.cswf.util;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading;
using System.Threading.Tasks;

[assembly: log4net.Config.XmlConfigurator(Watch = true)]
namespace io.vty.cswf.ffcm.console
{
    class Program
    {
        private static readonly ILog L = Log.New();
        static void Main(string[] args)
        {
            var conf = "conf/ffcm_c.properties";
            if (args.Length > 0)
            {
                conf = args[0];
            }
            var cfg = new FCfg();
            cfg.Load(conf, true);
            cfg.Print();
            var addr = cfg.Val("srv_addr", "");
            if (addr.Length < 1)
            {
                Console.WriteLine("the srv_addr is not setted");
                Environment.Exit(1);
                return;
            }
            L.I("starting ffcm...");

            //Samba.
            var lambah = new LambdaEvnH();
            var ffcm = new DocCov("FFCM", cfg, new SckDailer(addr).Dail, lambah);
            var ffcmh = new FFCM(ffcm, ffcm.Srv);
            ffcm.InitConfig();
            ffcm.StartMonitor();
            ffcm.StartWindowCloser();
            ffcm.Start();
            ffcm.StartProcSrv();
            var activated = false;
            if (cfg.Val("samba", "N") == "Y")
            {
                L.I("start initial samba...");
                var samba = Samba.AddVolume2(cfg.Val("samba_vol", ""), cfg.Val("samba_uri", ""),
                    cfg.Val("samba_user", ""), cfg.Val("samba_pwd", ""),
                    cfg.Val("samba_paths", ""));
                samba.Fail = (s, e) =>
                {
                    ffcm.ChangeStatus(DTM_C.DCS_UNACTIVATED);
                    activated = false;
                };
                samba.Success = (s) =>
                {
                    if (!activated)
                    {
                        ffcm.ChangeStatus(DTM_C.DCS_ACTIVATED);
                        activated = true;
                    }
                };
                new Thread(run_samba).Start();
            }
            else
            {
                activated = true;
            }
            lambah.OnLogin = (nr, token) =>
            {
                if (activated)
                {
                    ffcm.ChangeStatus(DTM_C.DCS_ACTIVATED);
                }
            };
            lambah.EndCon = (nr) =>
            {
                ffcm.ChangeStatus(DTM_C.DCS_UNACTIVATED);
            };
            var reboot = cfg.Val("reboot", "");
            if (reboot.Length > 0)
            {
                ProcKiller.Shared.OnHavingNotKill = (c) =>
                {
                    string output;
                    Exec.exec(out output, reboot);
                };
            }
            new Thread(run_hb).Start(ffcm);
            ffcm.Wait();
        }
        static void run_hb(object s)
        {
            var ffcm = (DocCov)s;
            while (true)
            {
                try
                {
                    var time = Util.Now();
                    ffcm.hb("DocCov");
                    time = Util.Now() - time;
                    if (time > 1000)
                    {
                        L.W("DocCov do hb success, {0}ms used", time);
                    }
                }
                catch (Exception e)
                {
                    L.W("{0}", e.Message);
                }
                Thread.Sleep(16000);
            }
        }
        static void run_samba(object s)
        {
            Samba.LoopChecker();
        }
    }
}
