# GoldDigger
![golddigger-github](https://github.com/Varppi/GoldDigger/assets/72181445/d47ceace-cb86-4c6c-b430-50880356b19e)

### GolDigger is an easy to use automated web crawler and directory bruteforcer.

## Installation
```
go install github.com/SpoofIMEI/Varppi@latest
```

#### Add ~/go/bin to your $PATH

## Usage examples
```
#Basic 
GoldDigger -u "https://somesite.com"

#Considers all URLs that have the string "sometimes.com" in it as part of the target scope
GoldDigger -u "https://somesite.com" -k somesite.com

#Saves website file URLs to an output file
GoldDigger -u "https://somesite.com" -q -o results

#Run directory bruteforce with a custom wordlist
GolDigger -u "https://somesite.com" -w ~/Seclists/Discovery/Web-Content/raft-medium-directories.txt
```

## Example output
<img src="https://github.com/Varppi/GoldDigger/assets/72181445/e5bc1f82-13a6-4964-a555-62a87f88b0ca" width=800></img>
