(function () {
    function initNavToggle(root) {
        const toggle = root.querySelector("[data-nav-toggle]");
        const menu = document.getElementById("nav-menu");

        if (!toggle || !menu || toggle.dataset.navReady === "true") {
            return;
        }

        toggle.dataset.navReady = "true";

        toggle.addEventListener("click", function () {
            const isOpen = toggle.getAttribute("aria-expanded") === "true";

            toggle.setAttribute("aria-expanded", String(!isOpen));
            menu.hidden = isOpen;
        });

        document.addEventListener("click", function (event) {
            if (menu.hidden) {
                return;
            }

            if (menu.contains(event.target) || toggle.contains(event.target)) {
                return;
            }

            toggle.setAttribute("aria-expanded", "false");
            menu.hidden = true;
        });

        document.addEventListener("keydown", function (event) {
            if (event.key !== "Escape" || menu.hidden) {
                return;
            }

            toggle.setAttribute("aria-expanded", "false");
            menu.hidden = true;
            toggle.focus();
        });

        menu.addEventListener("click", function (event) {
            const link = event.target.closest("a");

            if (!link) {
                return;
            }

            toggle.setAttribute("aria-expanded", "false");
            menu.hidden = true;
        });
    }

    document.addEventListener("DOMContentLoaded", function () {
        initNavToggle(document);
    });

    document.body.addEventListener("htmx:afterSwap", function (event) {
        initNavToggle(event.target);
    });
})();
