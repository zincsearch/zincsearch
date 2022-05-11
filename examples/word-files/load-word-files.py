import textract
import httpx
from os import walk

folder = "sample-files"

def main():
    zinc_server = "https://playground.dev.zincsearch.com/api/books/document"
    zinc_uid = "admin"
    zinc_pwd = "Complexpass#123"

    f = []
    for (dirpath, dirnames, filenames) in walk(folder):
        f.extend(filenames)
        break

    for file in f:
        print(file)
        text = textract.process(folder + "/" + file)
        text = text.decode("utf-8") 
        data = {
            "file": file,
            "text": text
        }

        httpx.put(zinc_server, json=data, auth=(zinc_uid, zinc_pwd))

if __name__ == "__main__":
    main()
