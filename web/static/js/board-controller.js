import { Chessground } from "/static/js/chessground.min.js";

(function () {
    function initBoard(el) {
        if (el.dataset.boardReady === "true") {
            return;
        }

        const ground = Chessground(el, {
            fen: el.dataset.fen || "start",
            orientation: el.dataset.orientation || "white",
            turnColor: el.dataset.turn || "white",
            movable: {
                free: false,
                color: el.dataset.turn || "white",
                dests: new Map(),
                events: {
                    after: function (orig, dest, metadata) {
                        submitMove(el, ground, orig, dest, metadata);
                    }
                }
            },
            highlight: {
                lastMove: true,
                check: true
            }
        });

        el._chessground = ground;
        el.dataset.boardReady = "true";
    }

    async function submitMove(el, ground, orig, dest, metadata) {
        const url = el.dataset.moveUrl;

        if (!url) {
            return;
        }

        const form = new FormData();
        form.set("orig", orig);
        form.set("dest", dest);

        if (metadata && metadata.promotion) {
            form.set("promotion", metadata.promotion);
        }

        const response = await fetch(url, {
            method: "POST",
            body: form,
            headers: {
                "HX-Request": "true"
            }
        });

        if (!response.ok) {
            ground.set({
                fen: el.dataset.fen
            });
            return;
        }

        const state = await response.json();

        applyState(el, ground, state);
    }

    function applyState(el, ground, state) {
        el.dataset.fen = state.fen || el.dataset.fen;
        el.dataset.turn = state.turn || el.dataset.turn;

        ground.set({
            fen: state.fen,
            turnColor: state.turn,
            lastMove: state.lastMove,
            check: state.check || false,
            movable: {
                color: state.movableColor || state.turn,
                dests: objectToDests(state.dests || {})
            }
        });
    }

    function objectToDests(obj) {
        const dests = new Map();

        Object.entries(obj).forEach(function ([from, squares]) {
            dests.set(from, squares);
        });

        return dests;
    }

    function initBoards(root) {
        root.querySelectorAll("[data-board]").forEach(initBoard);
    }

    document.addEventListener("DOMContentLoaded", function () {
        initBoards(document);
    });

    document.body.addEventListener("htmx:afterSwap", function (event) {
        initBoards(event.target);
    });
})();
