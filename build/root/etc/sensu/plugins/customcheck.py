#!/usr/bin/python
# make sure to install python pip and to pip install requests in /root/setup.sh
import requests
import optparse
import sys

def check_server(host, port):
        try:
            r = requests.get("http://"+host+":"+str(port)+"/ping")
            if r.status_code < 200 or r.status_code >= 300:
                print("server up, but returning invalid code %d" % r.status_code)
                return 1
            return 0
        except Exception as e:
            print(e)
            return 2


if __name__ == '__main__':
    parser = optparse.OptionParser()
    
    parser.add_option("--host", dest="host", type="string", default='localhost', help="host to check out", metavar="URL")
    parser.add_option("-p", "--port", dest="port", type="int", default='8080', help="port", metavar="PORT")
    
    (options, args) = parser.parse_args()
    check = check_server(options.host, options.port)
    
    sys.exit(check)
