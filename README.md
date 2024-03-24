# GoAddons

GoAddons is a command-line interface (CLI) application designed to enhance the World of Warcraft (WoW) experience by simplifying the management of game addons. Leveraging the power of a database, GoAddons facilitates the effortless discovery, update, and management of WoW addons.

## Features

- **Addon Management**: Easily list, search, add, or remove WoW addons to customize your gaming setup.
- **Updater Menu**: ~~Automatically checks and~~ updates your addons to the latest versions, ensuring you have the latest features and fixes.
- **About**: Learn more about GoAddons.

## Getting Started

To get started with GoAddons, follow these simple steps:

1. **Clone the repository**:

```bash
git clone https://github.com/HansenChristoffer/GoAddons.git
cd GoAddons
```

2. **Docker containers**:

Make sure you have Docker/Docker-engine and Docker-compose installed on your system. You will need to edit the docker-compose.yml file to fit your system.
The "kaasufouji-extract-volume" device path needs to be your systems path to your addons directory. Finally, change any [YOUR_HOST_NAME_HERE] to your actual user's name.

Now, run the following command:

```bash
docker-compose up -d
```

This command will pull the relevant Docker images, create volumes and create containers that GoAddons will need to use.

4. **Build the application**:

Make sure you have Go installed on your system. You can build GoAddons using the provided Makefile for convenience:

```bash
make release
```

This command compiles the application and places the binary in the `bin` directory.

3. **Run GoAddons**:

```bash
./bin/goaddons
```

Follow the on-screen prompts to manage your WoW addons.

## Note

If you're planning on doing any **Addon Management**, you need to make sure that your **'tanukiDB'** container is up and running.

Also, keep in mind that if you're going to run the **Updater**, you need to make sure the Docker container called **'goaddons_cdp'** is up and running.

You can check what containers are currently up-and-running by executing the following:

```bash
docker ps
```

## Usage

After starting GoAddons, you'll be presented with the main menu:

```
  »»» GoAddons «««

 1. Addon Management
 2. Updater Menu
 3. About
 X. Exit
```

Choose an option by entering the corresponding number or letter and pressing `Enter`.

### Addon Management

In the Addon Management menu, you can:

- List all addons
- Search for a specific addon
- Add a new addon
- Remove an existing addon

### Updater Menu

Start the updater ~~to check for~~ and apply updates to your addons.

### About

Displays information about GoAddons and its development.

## Contributing

We welcome contributions to GoAddons! If you have suggestions or encounter issues, please feel free to open an issue or submit a pull request.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details. Contributions and modifications are welcome, but redistributed versions must also be licensed under the GPL v3. While not legally required, giving credit where credit is due is highly encouraged as a sign of respect and appreciation for the contributors' efforts.

## Acknowledgments

A big thank you to the WoW gaming community and all contributors to the GoAddons project.
