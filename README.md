<div align="center">
  <h1>🐡 Easy Tunneler 🐡</h1>
  <p><i>The dead simple solution for securely accessing your stuff from anywhere in the world</i></p>
  <img width="304" alt="image" src="https://github.com/user-attachments/assets/58c6adf9-000c-4e31-86fe-fda2369ae40c">
</div>

---

Inspired by [wg-easy](https://github.com/wg-easy/wg-easy), Easy-Tunneler makes accessing local applications in another computer stupidly simple from *any* device. Specifically, it is a frontend for OpenSSH, utilizing its amazing local port forwarding feature to access locally hosted applications like Nextcloud and Immich remotely. It intends to be a simpler alternative to accessing local apps via VPNs.

# Features
- 🤖 Create and automatically manage SSH tunnels, automatically rebooting tunnels on timeouts
- 🍃 Clean, responsive, and performant UI created in Vue.js
- 🔒 Based on OpenSSH, the most widely used and audited SSH implementation in the world. Also does not require root privileges to function
- ⚡ Designed to be ran daily on phones; uses little to no resources!

# Building from Source
1. Clone this repository
`git clone https://github.com/yanghuaxuan/Easy-Tunneler/`
2. Install the latest version of [Go](https://go.dev/doc/) and [Node.js](https://nodejs.org/en)
3. Run the build script
`./build.sh`

## Running
After building from source:
1. cd into `build`
2. Run the `run.sh` script provided

# Usage
You must have OpenSSH installed, and also have key-based authentication setup for the server you're connecting to. Password-based authentication is not supported for security purposes.
- See this [guide](https://www.digitalocean.com/community/tutorials/how-to-configure-ssh-key-based-authentication-on-a-linux-server) on how to setup key based authentication

# Credits
- [The OpenBSD Pufferfish](https://www.openbsd.org/)
