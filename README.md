<div style="text-align:center"><img src="https://github.com/t3chn0m4g3/h0neytr4p/blob/main/logo.png?raw=true" /></div>


# Adjusted for T-Pot

### This fork of [h0neytr4p](https://github.com/pbssubhash/h0neytr4p) was revised for T-Pot with the following features:
- Add Docker support (Dockerfile, docker-compose.yml)
- Use a single logfile instead of two
- Log to JSON instead of CSV
- Enrich JSON log file with additional info (i.e. Cookies, Headers, Destination Port, etc.)
- Improve trap support on multiple / different ports
- Add payload handler (store payloads in payload folder) with pre-defined sizeLimit(s)

## Original [h0neytr4p](https://github.com/pbssubhash/h0neytr4p) work by

### Authors:
- @pbssubhash; [Twitter](https://twitter.com/pbssubhash) | [LinkedIn](https://in.linkedin.com/in/pbssubhash)

### Rule Contributors: 
- @me-godsky; [Twitter](https://twitter.com/me_godsky) | [LinkedIn](https://in.linkedin.com/in/aakashmadaan13)

# What is h0neytr4p? 

Honeytrap (a.k.a h0neytr4p) is an easy to configure, deploy honeypot for protecting against web recon and exploiting. 

# How does it work? 
Blue teams can create `trap` for each vulnerability or exploit or recon technique and place it in the `/traps` folder and restart h0neytr4p. This will automatically reload the configuration and start the h0neytr4p.

# What does it protect against?
h0neytr4p was primarly built to remove the pain of creating a vulnerable application for publicly facing honeypots. While there's no denying the fact that creating an end to end vulnerable application might have it's own advantages, we need something flexible, agile framework for trapping the notorious bad guys. Some of the common use-cases are:
- Let's say you received an advisory that some XXX group is targetting a web RCE 1day and you want to detect the exploitation or recon attempts, you are at the right place.
- You want to know who's scanning your external attack surface using the new cutting edge tools like nuclei or nmap? this tool got it covered.

## Run with docker 

```
git clone https://github.com/t3chn0m4g3/h0neytr4p
cd h0neytr4p
docker compose build
docker compose up

  /$$        /$$$$$$                                  /$$               /$$   /$$
 | $$       /$$$_  $$                                | $$              | $$  | $$
 | $$$$$$$ | $$$$\ $$ /$$$$$$$   /$$$$$$  /$$   /$$ /$$$$$$    /$$$$$$ | $$  | $$  /$$$$$$
 | $$__  $$| $$ $$ $$| $$__  $$ /$$__  $$| $$  | $$|_  $$_/   /$$__  $$| $$$$$$$$ /$$__  $$
 | $$  \ $$| $$\ $$$$| $$  \ $$| $$$$$$$$| $$  | $$  | $$    | $$  \__/|_____  $$| $$  \ $$
 | $$  | $$| $$ \ $$$| $$  | $$| $$_____/| $$  | $$  | $$ /$$| $$            | $$| $$  | $$
 | $$  | $$|  $$$$$$/| $$  | $$|  $$$$$$$|  $$$$$$$  |  $$$$/| $$            | $$| $$$$$$$/
 |__/  |__/ \______/ |__/  |__/ \_______/ \____  $$   \___/  |__/            |__/| $$____/
                                         /$$  | $$                               | $$
                                        |  $$$$$$/                [ v0.3 ]       | $$
                                         \______/                                |__/
 Built by a Red team, with <3
 Built by zer0p1k4chu & g0dsky (https://github.com/pbssubhash/h0neytr4p)
 Adjusted for T-Pot by t3chn0m4g3 (https://github.com/t3chn0m4g3/h0neytr4p)
 	
 [ Traps folder            ] -> [ traps/                        ]
 [ Logfile                 ] -> [ log/log.json                  ]
 [ Payloads folder         ] -> [ /opt/h0neytr4p/payloads/      ]
 [ Catch all payloads      ] -> [ false                         ]
 [ Payload multipart limit ] -> [ 103424                        ]
 [ Payload other limit     ] -> [ 11264                         ]

 Logging is configured and ready.
 Payload folder is configured and ready.
 [~>] Loaded 31 traps on Port:443. Let's get the ball rolling!
```
