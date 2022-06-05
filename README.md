Huna
===========

This repository holds a [web assembly](https://webassembly.org/) based [keepass](https://keepass.info/) client compatible with [uSuite](https://github.com/ChristianHering/uSuite).

It provides:

  * A WASM based KeePass client
  * A uSuite plugin example

Table of Contents:

  * [About](#about)
  * [Compiling from Source](#compiling-from-source)
  * [Contributing](#contributing)
  * [License](#license)

About
-----

Huna spawned from me not liking the keepass clients currently available for mobile linux platforms. However, I later found [KeeWeb](https://keeweb.info/) which is essentially what this project was going to be. Because of this, I'm shelving development of the project until future notice.

Compiling from Source
------------

In order to compile Huna from source, you'll need [uSuite](https://github.com/ChristianHering/uSuite) and then you'll need to add the following to uSuite's [main function](https://github.com/ChristianHering/uSuite/blob/master/main.go).
```Go
err := callbackTemplate(huna.Huna(mux, configuration.DataDir, callback))
	if err != nil {
		panic(err)
	}
```
Then, you'll need to compile the WASM part of Huna by going into the 'client' folder and running:
`GOOS=js GOARCH=wasm go build -o ./../asm/bin.wasm`
After that, you should be able to run `go run ./` from uSuite and have everything work as expected.

Contributing
------------

Contributions are always welcome. If you're interested in contributing, send me an email or submit a PR.

License
-------

This project is currently licensed under GPLv3. This means you may use our source for your own project, so long as it remains open source and is licensed under GPLv3.

Please refer to the [license](/LICENSE) file for more information.
