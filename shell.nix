{ system ? builtins.currentSystem
, sources ? import ./nix/sources.nix
}: let
	pkgs = import sources.nixpkgs {
		inherit system;
		config.packageOverrides = pkgs: {
			nur = import sources.NUR { inherit pkgs; };
			unstable = import sources.unstable { inherit pkgs; };
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
		pkgs.go-junit-report
		pkgs.go-task
		pkgs.goreleaser
		pkgs.grype
		pkgs.hadolint
		pkgs.jq
		pkgs.lefthook
		pkgs.markdownlint-cli
		pkgs.nur.repos.wwmoraes.codecov-cli-bin
		pkgs.oapi-codegen
		pkgs.svu
		pkgs.typos
		pkgs.unstable.container-structure-test
		pkgs.unstable.go
		pkgs.unstable.golangci-lint
		pkgs.unstable.hadolint-sarif
		pkgs.unstable.sarif-fmt
	] ++ lib.optionals (builtins.getEnv "CI" != "") [ # CI-only
	] ++ lib.optionals (builtins.getEnv "CI" == "") [ # local-only
		## TODO pkgsite
		## TODO https://github.com/xo/xo
		pkgs.cocogitto
		pkgs.curl
		pkgs.docker-client
		pkgs.eclint
		pkgs.flyctl
		pkgs.niv
		pkgs.nur.repos.wwmoraes.go-commitlint
		pkgs.nur.repos.wwmoraes.gopium
		pkgs.nur.repos.wwmoraes.goutline
		pkgs.nur.repos.wwmoraes.structurizr-cli
		pkgs.nur.repos.wwmoraes.structurizr-site-generatr
		pkgs.plantuml
		pkgs.redis
		pkgs.sqlc
		pkgs.unstable.delve
		pkgs.unstable.go-cover-treemap
		pkgs.unstable.golangci-lint-langserver
		pkgs.unstable.gopls
		pkgs.unstable.gotests
		pkgs.unstable.gotools
	];
}
