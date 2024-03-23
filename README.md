# GoAddons

GoAddons is a sophisticated command-line interface (CLI) application designed to enhance the World of Warcraft (WoW) gaming experience by simplifying the management of game addons. Leveraging the power of TanukiDB, GoAddons facilitates the effortless discovery, update, and management of WoW addons, ensuring gamers can focus on what truly matters: their gaming experience.

## Features

- **Addon Management**: Easily list, search, add, or remove WoW addons to customize your gaming setup.
- **Updater Menu**: Automatically checks and updates your addons to the latest versions, ensuring you have the latest features and fixes.
- **About**: Learn more about GoAddons, its philosophy, and the team behind it.

## Getting Started

To get started with GoAddons, follow these simple steps:

1. **Clone the repository**:

```bash
git clone https://github.com/yourusername/GoAddons.git
cd GoAddons
```

2. **Build the application**:

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

Start the updater to check for and apply updates to your addons.

### About

Displays information about GoAddons and its development.

## Contributing

We welcome contributions to GoAddons! If you have suggestions or encounter issues, please feel free to open an issue or submit a pull request.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details. Contributions and modifications are welcome, but redistributed versions must also be licensed under the GPL v3. While not legally required, giving credit where credit is due is highly encouraged as a sign of respect and appreciation for the contributors' efforts.

## Acknowledgments

A big thank you to the WoW gaming community and all contributors to the GoAddons project.