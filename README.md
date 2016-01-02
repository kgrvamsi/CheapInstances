# CheapInstances
CheatInstances is a tool to get the Spot instances from AWS through CLI

```
// With Options
Usage: cheapinstances -count="3" -instance="m3.medium" -price="0.10" -region="us-east-1" -zone="us-east-1c"
  -count string
    This is used for the number of instance requests you need
  -instance string
    This Represents the Instance Type
  -price string
    This Represents the Price to use for the Requested instance
   -region string
    This Represents the Region to Use
   -zone string
    This Represents the zone to use for the respective Region

// Without Options
./cheapinstances
Enter any one in the following:
	1) Get the price history
	2) Create the Spot Instance
	3) Check the Spot Instance Request History
	4) Cancel the Spot Instance Request
	5) Exit
```


## Folder Strucutre

```
├── Dockerfile
├── README.md
├── config.ini
├── glide.lock
├── glide.yaml
├── main.go
├── server
│   └── server.go
└── slack
    └── slack.go
```

## Third Party Integrations

Slack: We have integrated Slack Chat based integration on every request you make




