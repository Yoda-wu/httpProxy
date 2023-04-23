from selenium import webdriver
import time

options = webdriver.EdgeOptions()
options.use_chromium = True
# options.add_argument("--headless")
executable_path = r"D:/webdriver/msedgedriver.exe"

driver = webdriver.Edge(options=options, executable_path=executable_path)

js = """
    entries = window.performance.getEntries();
    console.log(entries)
    return entries
"""
driver.get("http://www.lib.scut.edu.cn/2016/1025/c8738a127507/page.htm")
entries = driver.execute_script(js)

for i in range(0, len(entries)):
    entry = entries[i]
    if entry["entryType"] == "paint":
        continue
    try:
        print("$" + str(i))
        print("name: " + entry["name"])
        print("size: %.2f" % (entry["transferSize"]) + "B")
        print("load time: %.2f" % (entry["duration"]) + "ms")
        print("#############################################")
    except:
        pass
time.sleep(10)
driver.close()