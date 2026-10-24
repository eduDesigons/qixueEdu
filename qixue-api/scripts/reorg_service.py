#!/usr/bin/env python3
"""Reorganize internal/service into admin/, webapp/, shared root.

Directory layout changes; package clause names stay the same so handler
packages (login/home/course/…) do not collide with service package names.
"""
from __future__ import annotations

import shutil
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SERVICE = ROOT / "internal" / "service"

# old_dir -> new_rel under service/  (package clause unchanged)
MOVES = {
    # admin*
    "adminbanner": "admin/banner",
    "admincategory": "admin/category",
    "admincert": "admin/cert",
    "adminchapter": "admin/chapter",
    "admincomment": "admin/comment",
    "admincommon": "admin/common",
    "admincourse": "admin/course",
    "adminfeedback": "admin/feedback",
    "adminhome": "admin/home",
    "adminlogin": "admin/login",
    "adminmedia": "admin/media",
    "adminmediacate": "admin/mediacate",
    "adminmember": "admin/member",
    "adminmessage": "admin/message",
    "adminnotice": "admin/notice",
    "adminops": "admin/ops",
    "adminpaper": "admin/paper",
    "adminpath": "admin/path",
    "adminrole": "admin/role",
    "adminstaff": "admin/staff",
    "adminverse": "admin/verse",
    # learner-facing
    "comment": "webapp/comment",
    "course": "webapp/course",
    "feedback": "webapp/feedback",
    "home": "webapp/home",
    "learn": "webapp/learn",
    "paper": "webapp/paper",
    "path": "webapp/path",
    "rank": "webapp/rank",
    "user": "webapp/user",
    "webauth": "webapp/auth",
}


def move_packages() -> None:
    for old, new_rel in MOVES.items():
        src = SERVICE / old
        dst = SERVICE / new_rel
        if not src.is_dir():
            if dst.is_dir():
                print(f"skip (already moved): {old}")
                continue
            raise SystemExit(f"missing source: {src}")
        if dst.exists():
            raise SystemExit(f"destination exists: {dst}")
        dst.parent.mkdir(parents=True, exist_ok=True)
        shutil.move(str(src), str(dst))
        print(f"moved {old} -> {new_rel}")


def rewrite_imports() -> None:
    # longest first
    old_names = sorted(MOVES.keys(), key=len, reverse=True)
    path_map = [
        (f"qixue-api/internal/service/{old}", f"qixue-api/internal/service/{MOVES[old]}")
        for old in old_names
    ]
    for path in ROOT.rglob("*.go"):
        text = path.read_text(encoding="utf-8")
        orig = text
        for old_path, new_path in path_map:
            text = text.replace(old_path, new_path)
        if text != orig:
            path.write_text(text, encoding="utf-8")
            print(f"imports: {path.relative_to(ROOT)}")


def main() -> None:
    move_packages()
    rewrite_imports()
    print("done")


if __name__ == "__main__":
    main()
