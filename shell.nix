let
  ## nix-prefetch-url --unpack <url>
  nixpkgs = fetchTarball {
    url = "https://github.com/NixOS/nixpkgs/archive/refs/tags/24.05.tar.gz";
    sha256 = "1lr1h35prqkd1mkmzriwlpvxcb34kmhc9dnr48gkm8hh089hifmx";
  };
  nixpkgs-unstable = fetchTarball {
    name = "nixos-unstable-a14c5d651cee9ed70f9cd9e83f323f1e531002db";
    # url = "https://github.com/NixOS/nixpkgs/archive/refs/heads/nixpkgs-unstable.tar.gz";
    url = "https://github.com/NixOS/nixpkgs/archive/a14c5d651cee9ed70f9cd9e83f323f1e531002db.tar.gz";
    sha256 = "1b2dwbqm5vdr7rmxbj5ngrxm7sj5r725rqy60vnlirbbwks6aahb";
  };
  kaizen-src = fetchTarball {
    name = "kaizen-01fb203b0905ed33b45a562a2cb5e5b6330044a8";
    url = "https://github.com/wwmoraes/kaizen/archive/01fb203b0905ed33b45a562a2cb5e5b6330044a8.tar.gz";
    sha256 = "04207l8g0p94jix2brwyhky1cscnd9w6vjn5dzzpfyv71wc2g0qa";
  };
  pkgs = import nixpkgs {};
  unstable = import nixpkgs-unstable {};
  kaizen = import kaizen-src { inherit pkgs; };
  inherit (pkgs) lib mkShell;
in mkShell {
  packages = with pkgs; [
    # TODO try commitlint-rs
    curl
    docker-client
    editorconfig-checker
    git
    go-task
    goreleaser
    grype
    hadolint
    jq
    kaizen.go-commitlint
    lefthook
    markdownlint-cli
    oapi-codegen
    python312Packages.codecov
    svu
    typos
    unstable.container-structure-test
    unstable.go
    unstable.golangci-lint
    unstable.hadolint-sarif
    unstable.sarif-fmt
    ## TODO github.com/wadey/gocovmerge
    ## TODO github.com/Khan/genqlient
    ## TODO github.com/xo/xo
  ] ++ lib.optionals (builtins.getEnv "CI" != "") [ # CI-only
  ] ++ lib.optionals (builtins.getEnv "CI" == "") [ # local-only
    # fish
    ## TODO pkgsite
    flyctl
    kaizen.structurizr-cli
    kaizen.structurizr-site-generatr
    plantuml
    redis
    unstable.gopls
  ];
}
