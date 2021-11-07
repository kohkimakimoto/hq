#!/usr/bin/env python
from __future__ import division, print_function, absolute_import, unicode_literals
try:
    from http.server import HTTPServer, SimpleHTTPRequestHandler
except ImportError:
    # fallback for python2
    from SimpleHTTPServer import SimpleHTTPRequestHandler
    from BaseHTTPServer import HTTPServer
import argparse, time

sleepTime = 0

# RequestHandler
class RequestHandler(SimpleHTTPRequestHandler, object):
    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()
        self.wfile.write(b"It works!")

    def do_POST(self):
        rawPostData = ""
        length = self.headers.get('Content-Length')
        if (length):
          rawPostData = self.rfile.read(int(length))
        jobId = self.headers.get('X-Hq-Job-Id')
        if not jobId:
          jobId = "unknown"

        print("--POST REQUEST BEGIN--")
        print("%s %s\n%s" % (self.command, self.path, self.headers))
        print("%s" % (rawPostData))
        print("--POST REQUEST END----")

        global sleepTime
        if sleepTime > 0:
            print("Sleeping for %d seconds" % sleepTime)
            time.sleep(sleepTime)

        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()
        self.wfile.write(b"Got a job: " + jobId)

# main
def main():
    parser = argparse.ArgumentParser(
        description="An example worker web app for HQ",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
An example worker web app for HQ.

"""
    )

    parser.add_argument('-H', '--host', dest="host", default="")
    parser.add_argument('-p', '--port', dest="port", type=int, default=8000)
    parser.add_argument('-s', '--sleep', dest="sleep", type=int, default=0)
    args = parser.parse_args()

    port = args.port
    host = args.host
    global sleepTime
    sleepTime = args.sleep

    httpd = HTTPServer((host, port), RequestHandler)
    print("Serving at port", port)
    httpd.serve_forever()

if __name__ == '__main__': main()
