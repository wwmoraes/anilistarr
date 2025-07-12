{
  pkgs,
  ...
}:
let
  inherit (pkgs) lib mkShell;
in
mkShell {
  packages =
    [
      pkgs.checkmake
      pkgs.docker-client
      pkgs.editorconfig-checker
      pkgs.git
      pkgs.go-junit-report
      pkgs.goEnv
      pkgs.goreleaser
      pkgs.grype
      pkgs.hadolint
      pkgs.jq
      pkgs.markdownlint-cli
      pkgs.moreutils
      pkgs.nur.repos.wwmoraes.codecov-cli-bin
      pkgs.nur.repos.wwmoraes.structurizr-cli
      pkgs.omnix
      pkgs.plantuml
      pkgs.ripgrep
      pkgs.semgrep
      pkgs.unstable.cocogitto
      pkgs.unstable.container-structure-test
      pkgs.unstable.go
      pkgs.unstable.go-cover-treemap
      pkgs.unstable.go-task
      pkgs.unstable.golangci-lint
      pkgs.unstable.hadolint-sarif
      pkgs.unstable.sarif-fmt
      pkgs.yq-go
    ]
    ++ lib.optionals (builtins.getEnv "CI" != "") [
      # CI-only
    ]
    ++ lib.optionals (builtins.getEnv "CI" == "") [
      # local-only
      # pkgs.anilistarr
      pkgs.curl
      pkgs.docker-client
      pkgs.eclint
      pkgs.flyctl
      pkgs.gomod2nix
      pkgs.nur.repos.wwmoraes.gopium
      pkgs.nur.repos.wwmoraes.goutline
      pkgs.nur.repos.wwmoraes.structurizr-site-generatr
      pkgs.omnix
      pkgs.plantuml
      pkgs.redis
      pkgs.unstable.delve
      pkgs.unstable.go-cover-treemap
      pkgs.unstable.gopls
      pkgs.unstable.gotests
      pkgs.unstable.gotools
    ];

  shellHook = ''
    if [ -n "$CI" ]; then
      export GOCACHE=$(go env GOCACHE)
      export GOMODCACHE=$(go env GOMODCACHE)
    fi

    cog install-hook --all --overwrite
  '';
}
