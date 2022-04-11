# Discord in Insomnia

## What is this?

This is a small project that auto-parses the [markdown discord documentation]() into [Insomnia](https://github.com/Kong/insomnia)'s formatting, including creating env variables for unified requests.

You can now test the discord api with a nice rest client ;)

## How to use

### Using pre-built json

`Output.json` is a pre-made import file. You can just import that, either by downloading it and importing it (preferences -> data -> import), or [through url]()

### Manually

1. Download & navigate into this project
2. `go get github.com/ShadiestGoat/discord-insomnia`
3. `go build`
4. Run the executable, for linux/mac it's `discord-insomnia`, for windows it has a .exe extention.
5. Import the output, `Output.json`, into the insomnia clinet (preferences -> data -> import)

You can have an optional step where you download the resources yourself. [They are on the official discord repo](https://github.com/discord/discord-api-docs/tree/main/docs/resources)

