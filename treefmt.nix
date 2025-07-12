{
  projectRootFile = "flake.nix";

  ## TODO custom: checkmake
  ## TODO custom: editorconfig-checker/eclint
  ## TODO custom: golangci-lint
  ## TODO custom: hadolint
  ## TODO custom: markdownlint  # beautysh
  # programs.fish_indent.enable = true;
  programs.gofmt.enable = true;
  programs.gofumpt.enable = true;
  programs.goimports.enable = true;
  programs.golines.enable = true;
  # programs.jsonfmt.enable = true;
  # keep-sorted
  # programs.mdformat.enable = true;
  programs.nixf-diagnose.enable = true;
  # oxipng
  # pinact
  # programs.shellcheck-posix.enable = true;
  # programs.shellcheck-bash.enable = true;
  # shfmt
  programs.nixfmt.enable = true;
  programs.statix.enable = true;
  # programs.taplo = {
  #   enable = true;
  #   excludes = [
  #     "gomod2nix.toml"
  #   ];
  # };
  programs.typos = {
    enable = true;
    excludes = [
      "**/*.gen.go"
      ".golangci.yaml"
      "CHANGELOG.md"
      "go.mod"
      "go.sum"
    ];
    configFile = builtins.toString ./.typos.toml;
  };
  programs.yamlfmt = {
    enable = true;
    excludes = [
      ".direnv"
      "tmp"
    ];
    settings = {
      gitignore_excludes = true;
      formatter = {
        type = "basic";
        indentless_arrays = true;
        retain_line_breaks_single = true;
        scan_folded_as_literal = true;
        trim_trailing_whitespace = true;
      };
    };
  };

}
