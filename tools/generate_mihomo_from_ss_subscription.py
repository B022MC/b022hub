#!/usr/bin/env python3
"""Generate a Mihomo config from a base64 ss:// subscription."""

from __future__ import annotations

import argparse
import base64
import sys
import urllib.parse
import urllib.request
from dataclasses import dataclass
from typing import Iterable


DEFAULT_TEST_URL = "http://www.gstatic.com/generate_204"
DEFAULT_INTERVAL = 600
DEFAULT_TIMEOUT = 20
INFO_NAME_KEYWORDS = ("到期", "剩余", "官网", "网址", "流量")


@dataclass(frozen=True)
class SSNode:
    name: str
    server: str
    port: int
    cipher: str
    password: str


def fetch_subscription(url: str, timeout: int) -> str:
    req = urllib.request.Request(
        url,
        headers={
            "User-Agent": "curl/8.0",
            "Accept": "*/*",
        },
    )
    with urllib.request.urlopen(req, timeout=timeout) as resp:
        return resp.read().decode("utf-8", "replace")


def decode_base64_text(raw: str) -> str:
    compact = "".join(raw.strip().splitlines())
    padding = (-len(compact)) % 4
    return base64.b64decode(compact + ("=" * padding)).decode("utf-8", "replace")


def decode_userinfo(encoded: str) -> tuple[str, str]:
    padding = (-len(encoded)) % 4
    decoded = base64.urlsafe_b64decode(encoded + ("=" * padding)).decode("utf-8")
    cipher, password = decoded.split(":", 1)
    return cipher, password


def sanitize_name(name: str, used: set[str]) -> str:
    cleaned = " ".join(urllib.parse.unquote(name).replace("+", " ").split()).strip()
    if not cleaned:
        cleaned = "node"
    candidate = cleaned
    suffix = 2
    while candidate in used:
        candidate = f"{cleaned} #{suffix}"
        suffix += 1
    used.add(candidate)
    return candidate


def is_info_name(name: str) -> bool:
    return any(keyword in name for keyword in INFO_NAME_KEYWORDS)


def parse_nodes(text: str) -> list[SSNode]:
    seen: set[tuple[str, int, str, str]] = set()
    used_names: set[str] = set()
    nodes: list[SSNode] = []

    for raw_line in text.splitlines():
        line = raw_line.strip()
        if not line.startswith("ss://"):
            continue
        parsed = urllib.parse.urlparse(line)
        if not parsed.username or not parsed.hostname or parsed.port is None:
            continue
        cipher, password = decode_userinfo(parsed.username)
        name = sanitize_name(parsed.fragment or "", used_names)
        if is_info_name(name):
            continue
        dedupe_key = (parsed.hostname, parsed.port, cipher, password)
        if dedupe_key in seen:
            continue
        seen.add(dedupe_key)
        nodes.append(
            SSNode(
                name=name,
                server=parsed.hostname,
                port=parsed.port,
                cipher=cipher,
                password=password,
            )
        )
    return nodes


def render_proxy(node: SSNode, indent: str = "  ") -> str:
    return "\n".join(
        [
            f"{indent}- name: {yaml_quote(node.name)}",
            f"{indent}  type: ss",
            f"{indent}  server: {yaml_quote(node.server)}",
            f"{indent}  port: {node.port}",
            f"{indent}  cipher: {yaml_quote(node.cipher)}",
            f"{indent}  password: {yaml_quote(node.password)}",
            f"{indent}  udp: true",
        ]
    )


def render_string_list(values: Iterable[str], indent: str) -> str:
    return "\n".join(f"{indent}- {yaml_quote(value)}" for value in values)


def yaml_quote(value: str) -> str:
    escaped = value.replace("\\", "\\\\").replace('"', '\\"')
    return f"\"{escaped}\""


def build_config(nodes: list[SSNode], test_url: str, interval: int) -> str:
    if not nodes:
        raise ValueError("subscription contains no usable ss nodes")

    names = [node.name for node in nodes]
    parts = [
        "mixed-port: 7890",
        "allow-lan: true",
        "bind-address: '*'",
        "mode: rule",
        "log-level: info",
        "external-controller: 0.0.0.0:9090",
        "proxies:",
        "\n".join(render_proxy(node) for node in nodes),
        "proxy-groups:",
        "  - name: \"Auto\"",
        "    type: url-test",
        "    url: " + yaml_quote(test_url),
        f"    interval: {interval}",
        "    proxies:",
        render_string_list(names, "      "),
        "  - name: \"SubscriptionPool\"",
        "    type: load-balance",
        "    strategy: round-robin",
        "    url: " + yaml_quote(test_url),
        f"    interval: {interval}",
        "    proxies:",
        render_string_list(names, "      "),
        "rules:",
        "  - MATCH,SubscriptionPool",
        "",
    ]
    return "\n".join(parts)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--url", required=True, help="Base64 ss subscription URL")
    parser.add_argument("--output", required=True, help="Path to write Mihomo YAML config")
    parser.add_argument("--test-url", default=DEFAULT_TEST_URL, help="Health-check URL")
    parser.add_argument(
        "--interval",
        type=int,
        default=DEFAULT_INTERVAL,
        help="Health-check interval in seconds",
    )
    parser.add_argument(
        "--timeout",
        type=int,
        default=DEFAULT_TIMEOUT,
        help="HTTP fetch timeout in seconds",
    )
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    raw = fetch_subscription(args.url, timeout=args.timeout)
    decoded = decode_base64_text(raw)
    nodes = parse_nodes(decoded)
    config = build_config(nodes, test_url=args.test_url, interval=args.interval)
    with open(args.output, "w", encoding="utf-8") as fh:
        fh.write(config)
    print(f"wrote {len(nodes)} nodes to {args.output}")
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as exc:  # pragma: no cover
        print(f"error: {exc}", file=sys.stderr)
        raise SystemExit(1)
