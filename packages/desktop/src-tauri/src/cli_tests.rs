use super::*;

#[test]
fn parses_clean_json_stdout() {
    let parsed: serde_json::Value = parse_cli_output(
        r#"[{"id":"BLAME-001","title":"BLAME-001: demo"}]"#,
        "",
        Some(0),
    )
    .expect("valid JSON stdout should parse");

    assert_eq!(parsed[0]["id"], "BLAME-001");
}

#[test]
fn rejects_stderr_even_when_exit_code_is_zero() {
    let error = parse_cli_output::<serde_json::Value>(
        "",
        "Error: unsupported path: GET /blames",
        Some(0),
    )
    .expect_err("stderr means the CLI result is not successful");

    assert!(error.contains("CLI wrote to stderr"));
}
