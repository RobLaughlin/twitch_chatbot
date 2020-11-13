# Twitch Chatbot
An IRC bot for Twitch!

## Generate Access/Refresh Tokens
To connect your bot account to twitch, you first need to generate access/refresh tokens. Running the following command below will create an auth.json file with all the necessary configuration required to run the bot. Make sure you replace the CLIENT_ID and CLIENT_SECRET fields in the auth.env before starting the authentication process. 
```
cd auth
go run auth.go
```

## Connect to Twitch
After the intial authentication setup, connecting to Twitch is fairly easy. After the proper environment variables are set for USER and PASS (being the oauth authentication token), run the following command:
```
go run bot.go
```