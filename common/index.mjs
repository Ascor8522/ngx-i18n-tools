#!/usr/bin/env -S node --disable-warning=ExperimentalWarning

"use strict";

import {WASI} from "node:wasi";
import {argv, env} from "node:process";
import {open} from "node:fs/promises";

const wasi = new WASI({
    version: "preview1",
    args: argv,
    env,
    preopens: {
        "/": process.cwd(),
    },
});

await Promise
    .resolve(new URL("./TO_REPLACE_TOOL.wasm", import.meta.url))
    .then(open)
    .then(fd => fd.readableWebStream({autoClose: true}))
    .then(rs => new Response(rs, {headers: {"Content-Type": "application/wasm"}}))
    .then(res => WebAssembly.instantiateStreaming(res, wasi.getImportObject()))
    .then(({instance}) => wasi.start(instance));
