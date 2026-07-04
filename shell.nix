{
  pkgs,
  ...
}:
rec {
  default = pkgs.mkShell {
    nativeBuildInputs = [
      # keep-sorted start
      (pkgs.mkGoEnv { pwd = ./.; })
      pkgs.docker-client
      pkgs.eclint
      pkgs.editorconfig-checker
      pkgs.git
      pkgs.goreleaser
      pkgs.gotestdox
      pkgs.grype
      pkgs.hadolint
      pkgs.jq
      pkgs.moreutils
      pkgs.remake
      pkgs.ripgrep
      pkgs.semgrep
      pkgs.unstable.container-structure-test
      pkgs.unstable.go
      pkgs.unstable.golangci-lint
      pkgs.valkey
      pkgs.unstable.sarif-fmt
      # keep-sorted end
    ];
  };

  ci = default.overrideAttrs (
    final: prev: {
      nativeBuildInputs = [
        # keep-sorted start
        pkgs.go-junit-report
        pkgs.nur.repos.wwmoraes.codecov-cli-bin
        pkgs.unstable.hadolint-sarif
        # keep-sorted end
      ]
      ++ prev.nativeBuildInputs;

      shellHook = ''
        export GOCACHE=$(go env GOCACHE)
        export GOMODCACHE=$(go env GOMODCACHE)
      '';
    }
  );

  terminal = default.overrideAttrs (
    final: prev: {
      nativeBuildInputs = [
        # pkgs.anilistarr
        # keep-sorted start
        pkgs.curl
        pkgs.gomod2nix
        # pkgs.nur.repos.wwmoraes.gopium
        # pkgs.nur.repos.wwmoraes.goutline
        pkgs.nur.repos.wwmoraes.structurizr-cli
        pkgs.nur.repos.wwmoraes.structurizr-site-generatr
        pkgs.plantuml
        pkgs.unstable.cocogitto
        pkgs.unstable.go-cover-treemap
        pkgs.unstable.gotests
        pkgs.unstable.gotools
        # keep-sorted end
      ]
      ++ prev.nativeBuildInputs;

      shellHook = ''
        ./configure
      '';
    }
  );
}
