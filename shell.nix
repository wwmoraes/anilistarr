let
  pkgs = import (fetchTarball {
    url = "https://github.com/NixOS/nixpkgs/archive/refs/tags/24.05.tar.gz";
    sha256 = "1lr1h35prqkd1mkmzriwlpvxcb34kmhc9dnr48gkm8hh089hifmx";
  }) {};
  unstable = import (fetchTarball {
    name = "nixos-unstable-a14c5d651cee9ed70f9cd9e83f323f1e531002db";
    url = "https://github.com/NixOS/nixpkgs/archive/a14c5d651cee9ed70f9cd9e83f323f1e531002db.tar.gz";
    sha256 = "1b2dwbqm5vdr7rmxbj5ngrxm7sj5r725rqy60vnlirbbwks6aahb";
  }) {};
  kaizen = import (fetchTarball {
    name = "kaizen-8075b45edf93d8f95a00958fd3a1cc606ba3405c";
    url = "https://github.com/wwmoraes/kaizen/archive/8075b45edf93d8f95a00958fd3a1cc606ba3405c.tar.gz";
    sha256 = "1w6fkd5kqm2l2aij31kd9ddk84phbcir7flrbha6mxw2y4j36z09";
  }) { inherit pkgs; };
  inherit (pkgs) lib mkShell;
in mkShell {
  packages = [
    ## TODO github.com/Khan/genqlient
    ## TODO github.com/wadey/gocovmerge
    ## TODO github.com/xo/xo
    ## TODO try commitlint-rs
    kaizen.go-commitlint
    pkgs.curl
    pkgs.docker-client
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
    pkgs.python312Packages.codecov
    pkgs.svu
    pkgs.typos
    unstable.container-structure-test
    unstable.go
    unstable.gotests
    unstable.golangci-lint
    unstable.hadolint-sarif
    unstable.sarif-fmt
  ] ++ lib.optionals (builtins.getEnv "CI" != "") [ # CI-only
  ] ++ lib.optionals (builtins.getEnv "CI" == "") [ # local-only
    ## TODO pkgsite
    kaizen.gopium
    kaizen.goutline
    kaizen.structurizr-cli
    kaizen.structurizr-site-generatr
    pkgs.flyctl
    pkgs.plantuml
    pkgs.redis
    unstable.gopls
  ];
}
