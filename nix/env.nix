let
	pkgs = import <nixpkgs> {
		config = {};
		overlays = [];
	};
in
pkgs.mkShell {
	packages = with pkgs; [
		bash
		zsh
		procps
		tree
		git
		nano
		curl
		gzip
		zopfli
		brotli
	];
	GIT_EDITOR = "${pkgs.nano}/bin/nano";
	shellHook = ''
		export ZDOTDIR=$PWD/nix/zsh
		export XDG_CONFIG_HOME=$PWD/nix
		export PATH=$PWD:$PATH
	'';
}