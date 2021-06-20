<div style="text-align:center"><img src="https://github.com/pbssubhash/h0neytr4p/blob/main/logo.png?raw=true" /></div>

## What is h0neytr4p? 
Honeytrap (a.k.a h0neytr4p) is an easy to configure, deploy honeypot for protecting against web recon and exploiting. 

## How does it work? 
Blue teams can create `trap` for each vulnerability or exploit or recon technique and place it in the `/traps` folder and restart h0neytr4p. This will automatically reload the configuration and start the h0neytr4p.

## What does it protect against?
h0neytr4p was primarly built to remove the pain of creating a vulnerable application for publicly facing honeypots. While there's no denying the fact that creating an end to end vulnerable application might have it's own advantages, we need something flexible, agile framework for trapping the notorious bad guys. Some of the common use-cases are:
- Let's say you received an advisory that some XXX group is targetting a web RCE 1day and you want to detect the exploitation or recon attempts, you are at the right place.
- You want to know who's scanning your external attack surface using the new cutting edge tools like nuclei or nmap? this tool got it covered.


## How to deploy it? 
The tool was build on top of Golang which means it can be easily compiled to your server/machine platform and architecture. 

#### To Build: 

```
go build main.go
./main -h

 /$$        /$$$$$$                                  /$$               /$$   /$$
| $$       /$$$_  $$                                | $$              | $$  | $$
| $$$$$$$ | $$$$\ $$ /$$$$$$$   /$$$$$$  /$$   /$$ /$$$$$$    /$$$$$$ | $$  | $$  /$$$$$$
| $$__  $$| $$ $$ $$| $$__  $$ /$$__  $$| $$  | $$|_  $$_/   /$$__  $$| $$$$$$$$ /$$__  $$
| $$  \ $$| $$\ $$$$| $$  \ $$| $$$$$$$$| $$  | $$  | $$    | $$  \__/|_____  $$| $$  \ $$
| $$  | $$| $$ \ $$$| $$  | $$| $$_____/| $$  | $$  | $$ /$$| $$            | $$| $$  | $$
| $$  | $$|  $$$$$$/| $$  | $$|  $$$$$$$|  $$$$$$$  |  $$$$/| $$            | $$| $$$$$$$/
|__/  |__/ \______/ |__/  |__/ \_______/ \____  $$   \___/  |__/            |__/| $$____/
                                         /$$  | $$                              | $$
       Built by a Red team, with <3     |  $$$$$$/                              | $$
             h0neytr4p v0.1             \______/                               |__/
        Built by zer0p1k4chu & g0dsky
    https://github.com/pbssubhash/h0neyt4p

Wrong Arguments.. Exiting Now
  -help string
        Print Help (default "Print Help")
  -log string
        Log file - It's a string. (default "Default")
  -output string
        Output file - It's a string. (default "Default")
  -traps string
        Traps folder - It's a string. (default "Default")
  -verbose string
        Use -verbose=false for disabling streaming output; by default it's true (default "true")
```

