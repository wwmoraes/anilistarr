{
  pkgs,
  ...
}:
rec {
  default = pkgs.mkShell {
    nativeBuildInputs = [
      (pkgs.mkGoEnv { pwd = ./.; })
      # keep-sorted start
      pkgs.git
      pkgs.gomod2nix
      pkgs.goreleaser
      pkgs.gotestdox
      pkgs.grype
      pkgs.hadolint
      pkgs.jq
      pkgs.moreutils
      pkgs.remake
      pkgs.ripgrep
      pkgs.semgrep
      pkgs.unstable.cocogitto
      pkgs.unstable.container-structure-test
      pkgs.unstable.go
      pkgs.unstable.golangci-lint
      pkgs.unstable.hadolint-sarif
      pkgs.unstable.sarif-fmt
      pkgs.valkey
      # keep-sorted end
    ];
  };

  ci = default.overrideAttrs (
    final: prev: {
      nativeBuildInputs = [
        # keep-sorted start
        pkgs.docker-client
        pkgs.go-junit-report
        pkgs.nur.repos.wwmoraes.codecov-cli-bin
        pkgs.qemu
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
        # keep-sorted start
        pkgs.curl
        pkgs.nur.repos.wwmoraes.structurizr-cli
        pkgs.nur.repos.wwmoraes.structurizr-site-generatr
        pkgs.plantuml
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
