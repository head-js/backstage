use super::*;

#[test]
fn routes_get_blames_to_list_blames_action() {
    let action = route_edge_action("get", "/blames").expect("GET /blames should be supported");

    assert_eq!(action, EdgeAction::ListBlames);
}

#[test]
fn rejects_unknown_edge_route() {
    let error = route_edge_action("POST", "/blames").expect_err("POST /blames is unsupported");

    assert!(error.contains("unsupported edge route: POST /blames"));
}

#[test]
fn maps_gitea_blame_to_edge_blame() {
    let parsed = transform_list_blames_output(
        r####"[{"id":"BLAME-005","name":"Needs decision","gitea":{"owner":"mall-view-consumer","repo":"PLAN-003","body":"### Ask\nPick one","updatedAt":"2026-06-08T11:15:46+08:00"}}]"####,
        "",
        Some(0),
    )
    .expect("valid blame output should map to edge DTO");

    assert_eq!(parsed[0]["id"], "mall-view-consumer-PLAN-003-BLAME-005");
    assert_eq!(parsed[0]["blameId"], "BLAME-005");
    assert_eq!(parsed[0]["name"], "Needs decision");
    assert_eq!(parsed[0]["appId"], "mall-view-consumer");
    assert_eq!(parsed[0]["planId"], "PLAN-003");
    assert_eq!(parsed[0]["context"], "### Ask\nPick one");
    assert_eq!(parsed[0]["updatedAt"], "06-08 11:15:46");
    assert!(parsed[0].get("gitea").is_none());
}

#[test]
fn rejects_gitea_blame_without_required_id() {
    let error = transform_list_blames_output(
        r####"[{"name":"Needs decision","gitea":{"owner":"mall-view-consumer","repo":"PLAN-003","body":"### Ask\nPick one","updatedAt":"2026-06-08T11:15:46+08:00"}}]"####,
        "",
        Some(0),
    )
    .expect_err("missing required blame id should fail parsing");

    assert!(error.contains("missing field `id`"));
}
