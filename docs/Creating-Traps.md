## What is a trap?

A trap is a `.json` file which consists of behaviour related a trap. For example, you can define that it should trap any request that matches `/jenkins`. All requests to `/jenkins` will be captured.

## Structure of a trap
Each trap must contain the following details in the JSON format. Most of the elements are self explanatory. A summarised version of all supported elements are mentioned below.

|S.No | Name of the attribute  | Description of the attribute  | Required/Optional  |  Data Type |
|---|---|---|---|---|
| 1.  |  BasicInfo | This is a nested json containing the basic information of the trap  | Required  |  JSON object |
| 2.  |  BasicInfo.Name | This is present inside basic info wrapper. This contains the name of the trap. Has to be unique across all the traps for better detection  |  Required | String  |
| 3.  | BasicInfo.Port  | The port where the trap listener should be started.  | Required   | String  |
| 4.  | BasicInfo.Protocol  |  The Protocol for the trap. Currently supported: HTTP  |  Required |  String |
| 4.  | BasicInfo.MitreAttackTags  |  For now, this will have external facing web compromise technique ID. However, it's reserved for future traps which will support more protocols |  Required | Strings seperated with a comma  |
| 5.  | BasicInfo.References  | Any reference URL for the attack for better detection and analytics.  | Required  | String  |
| 6.  |  BasicInfo.Description | Description of the attack which is being trapped.  | Required  | String  |
| 7.  |  BasicInfo.RiskRating | RiskRating for the trap; No defined values yet but generally accepted: Critical, High, Medium, Low, Info  | Required  | String  |
| 8.  |  Behaviour | It's a json object containing a pair of request and response.  | Required  | JSON array  |
| 9.  | Behaviour.Request  |  It's a json object containing the required request behaviour | Required  | JSON Object  |
| 10..  |  Beahviour.Request.Url |  URL Path.  You can define static or wildcards. [See reference below for how-to] | Required  | String  |
| 11.  |  Behaviour.Request.Method | Method; Supported Values: GET, POST, DELETE, PUT  |  Required |  String |
| 12.  | Behaviour.Request.Headers  | Request Headers. You can define static or wildcards. [See reference below for how-to]; Use an empty `{}` for empty headers | Required  |  JSON Object |
| 13.  |  Beahviour.Request.Params |  It's a key value pair json object. You can define the parameters irrespective of GET/POST and Content-Types. Backend handles it automatically. Use an empty `{}` for empty parameters. | Required  | JSON Object  |
| 14.  |  Beahviour.Response |  It's a json object containing the required response behaviour | Required  | JSON Object  |
| 15.  |  Beahviour.Request.StatusCode | Status code for response;  |  Required | String  |
| 16.  |  Beahviour.Request.Body | File content or location of the file;   |   |   |
| 17.  |  Beahviour.Request.Type | file or string  | Required | String  |
| 18.  |  Beahviour.trap | "true" if you want to trap, "false" if you don't want to   | Required  | String  |


#### Example Trap:
```
{
    "BasicInfo": {
        "Name":"jenkins_home",
        "Port":"443",
        "Protocol":"HTTP",
        "MitreAttackTags":"",
        "References":"",
        "RiskRating":"Critical",
        "Description":""
    },
    "Behaviour":
    [
        {
            "Request": {
                "Url":"/jenkins*",
                "Method": "GET",
                "Headers":{"User-Agent":"*"},
                "Params":{}
            },
            "Response": {
                "StatusCode": 302,
                "Body": "traps/assets/jenkins/default.html",
                "Type":"file"
            },
            "trap": "true"
        }
    ]
}
```

### Defining pattern inside an attribute:

You might've observed, `$` and `^` inside the trap, that's because h0neytr4p uses golang's regex for parsing your trap. 

###### Quick Walkthrough:

`.*` - wildcard 
`^` - starting of the string
`$` - ending of the string

Basically, 
- Let's say you want to match `/jenkins` in the Url field, you will use `/jenkins`. You can use `*` for defining a wildcard entry.
- Let's say you want to match `/wp-admin/login` in the Url field, you will use `/wp-admin/login`. 
- Let's say you want to match `/login.php?id=1'` and `/login.php?id=<ANY_Number>'`, you can use `^/login.php?id=*'` as the pattern.

The same goes for headers and params.

More examples: 

- You want to create a header list which accepts anything that starts with mozilla.: `"Headers": {"User-Agent":"Mozilla*"}` will be your header value.
- You want to create a parameter set which accepts username: anything starting with admin and password: password only.: `"Params":{"username":"*admin*","password":"password"}`


### Isn't there an automated way of converting requests into traps? 
Coming soon.

### Is there a way to verify if my trap syntax is right?
You can use any JSON lint tool. Specific syntax based checking is coming soon.
