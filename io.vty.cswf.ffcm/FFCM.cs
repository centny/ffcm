using io.vty.cswf.netw;
using io.vty.cswf.netw.dtm;
using io.vty.cswf.netw.http;
using io.vty.cswf.util;
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Web;
using io.vty.cswf.netw.r;
using io.vty.cswf.doc;

namespace io.vty.cswf.ffcm
{
    public class FFCM
    {
        public DTM_C DTMC { get; protected set; }
        public Server Srv { get; protected set; }
        public FFCM(DTM_C dtmc, Server srv)
        {
            this.DTMC = dtmc;
            this.Srv = srv;
            this.Srv.AddH("^/notify(\\?.*)?", this.OnFfProc);
            this.Srv.AddH("^/status(\\?.*)?", this.OnStatus);
        }

        public virtual HResult OnStatus(Request r)
        {
            r.res.StatusCode = 200;
            r.res.ContentEncoding = Encoding.UTF8;
            r.res.ContentType = "application/json";
            var data = new Dict();
            var tasks = new Dict();
            tasks.Add("queued", TaskPool.Shared.Queued.Count);
            tasks.Add("running", TaskPool.Shared.Running.Count);
            data.Add("tasks", tasks.data);
            data.Add("word_c", WordCov.Cached.Count);
            r.WriteLine(Json.stringify(data.data));
            r.Flush();
            return HResult.HRES_RETURN;
        }

        public virtual HResult OnFfProc(Request r)
        {
            var args = HttpUtility.ParseQueryString(r.req.Url.Query);
            var tid = args.Get("tid");
            var duration_ = args.Get("duration");
            if (String.IsNullOrWhiteSpace(tid) || String.IsNullOrWhiteSpace(duration_))
            {
                r.res.StatusCode = 400;
                r.WriteLine("the tid/duration is required");
                return HResult.HRES_RETURN;
            }
            float duration = 0;
            if (!float.TryParse(duration_, out duration))
            {
                r.res.StatusCode = 400;
                r.WriteLine("the duration must be float");
                return HResult.HRES_RETURN;
            }
            StreamReader reader = new StreamReader(r.req.InputStream);
            String line = null;
            var frame = new Dict();
            while ((line = reader.ReadLine()) != null)
            {
                line = line.Trim();
                if (line.Length < 1)
                {
                    continue;
                }
                var kvs = line.Split(new char[] { '=' }, 2);
                var key = kvs[0].Trim();
                if (kvs.Length < 2)
                {
                    frame[key] = "";
                }
                else
                {
                    frame[key] = kvs[1].Trim();
                }
                if (key != "progress")
                {
                    continue;
                }
                var ms = frame.Val<float>("out_time_ms", 0);
                this.DTMC.NotifyProc(tid, ms / duration);
            }
            r.res.StatusCode = 200;
            r.WriteLine("OK");
            r.Flush();
            return HResult.HRES_RETURN;
        }
    }
}
