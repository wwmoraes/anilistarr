{ system ? builtins.currentSystem
}: let
	pkgs = import (fetchTarball {
		url = "https://github.com/NixOS/nixpkgs/archive/refs/tags/24.05.tar.gz";
		sha256 = "1lr1h35prqkd1mkmzriwlpvxcb34kmhc9dnr48gkm8hh089hifmx";
	}) {
		inherit system;
		config.packageOverrides = pkgs: {
			nur = import (builtins.fetchTarball {
				url = "https://github.com/nix-community/NUR/archive/master.tar.gz";
				# sha256 = "0s9is965rv1knq17axd9s1y4l4h81d7dw6s0zsy5j6qwyb0kh703";
			}) { inherit pkgs; };
			unstable = import (fetchTarball {
				name = "nixos-unstable-a14c5d651cee9ed70f9cd9e83f323f1e531002db";
				url = "https://github.com/NixOS/nixpkgs/archive/a14c5d651cee9ed70f9cd9e83f323f1e531002db.tar.gz";
				sha256 = "1b2dwbqm5vdr7rmxbj5ngrxm7sj5r725rqy60vnlirbbwks6aahb";
			}) { inherit system pkgs; };
		};
	};
	inherit (pkgs) lib mkShell;
in mkShell {
	packages = [
	## TODO github.com/Khan/genqlient
	## TODO github.com/wadey/gocovmerge
	## TODO try commitlint-rs
	pkgs.editorconfig-checker
	pkgs.git
	pkgs.go-task
	pkgs.goreleaser
	pkgs.grype
	pkgs.hadolint
	pkgs.jq
	pkgs.lefthook
	pkgs.markdownlint-cli
	pkgs.oapi-codegen
	pkgs.svu
	pkgs.typos
	pkgs.unstable.container-structure-test
	pkgs.unstable.go
	pkgs.unstable.go-junit-report
	pkgs.unstable.golangci-lint
	pkgs.unstable.hadolint-sarif
	pkgs.unstable.sarif-fmt
	pkgs.nur.repos.wwmoraes.codecov-cli-bin
	] ++ lib.optionals (builtins.getEnv "CI" != "") [ # CI-only
	] ++ lib.optionals (builtins.getEnv "CI" == "") [ # local-only
	## TODO pkgsite
	## TODO https://github.com/xo/xo
	pkgs.eclint
	pkgs.curl
	pkgs.docker-client
	pkgs.flyctl
	pkgs.nur.repos.wwmoraes.go-commitlint
	pkgs.nur.repos.wwmoraes.gopium
	pkgs.nur.repos.wwmoraes.goutline
	pkgs.nur.repos.wwmoraes.structurizr-cli
	pkgs.nur.repos.wwmoraes.structurizr-site-generatr
	pkgs.plantuml
	pkgs.redis
	pkgs.unstable.delve
	pkgs.unstable.golangci-lint-langserver
	pkgs.unstable.gopls
	pkgs.unstable.gotests
	pkgs.unstable.gotools
	];
}
