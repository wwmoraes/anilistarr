{
  mkFormatterModule,
  pkgs,
  ...
}:
{
  imports = [
    (mkFormatterModule {
      name = "checkmake";
      package = "checkmake";
      args = [ ];
      includes = [
        "Makefile"
      ];
    })
    (mkFormatterModule {
      name = "hadolint";
      package = "hadolint";
      args = [ ];
      ## TODO generate config file
      includes = [
        "Dockerfile"
        "*.Dockerfile"
        "Dockerfile.*"
      ];
    })
  ];

  projectRootFile = "flake.nix";

  programs.checkmake.enable = true;
  ## TODO custom: editorconfig-checker/eclint
  ## TODO custom: golangci-lint
  programs.gofmt.enable = true;
  programs.gofumpt.enable = true;
  programs.goimports.enable = true;
  programs.golines.enable = true;
  programs.hadolint.enable = true;
  programs.keep-sorted.enable = true;
  programs.mdformat = {
    enable = true;
    package = pkgs.mdformat.withPlugins (
      ps: with ps; [
        mdformat-frontmatter
      ]
    );
    settings = {
      number = true;
      wrap = 80;
    };
  };
  programs.nixf-diagnose.enable = true;
  # oxipng
  # pinact
  programs.nixfmt.enable = true;
  programs.statix.enable = true;
  # programs.taplo = {
  #   enable = true;
  #   excludes = [
  #     "gomod2nix.toml"
  #   ];
  #   package = pkgs.unstable.taplo;
  #   settings = {
  #     formatting = {
  #       reorder_keys = true;
  #     };
  #   };
  # };
  programs.typos = {
    enable = true;
    excludes = [
      ".golangci.yaml"
      "CHANGELOG.md"
      "go.mod"
      "go.sum"
    ];
    configFile = builtins.toString ./.typos.toml;
  };
  programs.yamlfmt = {
    enable = true;
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
  settings.global.excludes = [
    ".direnv"
    "**/*.gen.go"
    "internal/drivers/sqlite/model/db.go"
    "internal/drivers/sqlite/model/models.go"
    "internal/drivers/sqlite/model/queries.sql.go"
    "tmp"
  ];
}
