using io.vty.cswf.log;
using io.vty.cswf.netw.sck;
using io.vty.cswf.util;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

[assembly: log4net.Config.XmlConfigurator(Watch = true)]
namespace io.vty.cswf.ffcm.console
{
    class Program
    {
        static void Main(string[] args)
        {
            var conf = "conf/ffcm_c.properties";
            if (args.Length > 0)
            {
                conf = args[0];
            }
            var cfg = new FCfg();
            cfg.Load(conf, true);
            Console.WriteLine(cfg);
            var addr = cfg.Val("srv_addr", "");
            if (addr.Length < 1)
            {
                Console.WriteLine("the srv_addr is not setted");
                Environment.Exit(1);
                return;
            }
            ILog L = Log.New();
            L.I("starting ffcm...");
            var ffcm = new FFCM("FFCM", cfg, new SckDailer(addr).Dail);
            ffcm.Start();
            ffcm.StartProcSrv();
            ffcm.Wait();
        }
    }
}
