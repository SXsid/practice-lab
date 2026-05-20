// ─── CONFIG ───────────────────────────────────────────────────────────────────
const FIREBASE_CONFIG = {
  apiKey: "AIzaSyBWNuaDcXkFi0b37acxgkBlVpEPpcudW0Y",
  authDomain: "fcm101-5e35e.firebaseapp.com",
  projectId: "fcm101-5e35e",
  storageBucket: "fcm101-5e35e.firebasestorage.app",
  messagingSenderId: "106012606884",
  appId: "1:106012606884:web:4bf44c9a1cb69c75ed8131",
};

// Firebase Console → Project Settings → Cloud Messaging → Web Push certificates
const VAPID_KEY =
  "BIMGg3ErMyrfFjEqq4BtkV92ublkB-Lm-EDiPx3AoYoWOSWGinaAURQ02f1EmOPl0SJg9B6E2cT0-oXGb4WmHBo";

// Optional: your backend endpoint to save {user_id, token}
// Set to null to skip saving (token still works for FCM)
const BACKEND_SAVE_URL = "http://localhost:8080/api/v1"; // e.g. "https://yourapi.com/save-token"
// ─────────────────────────────────────────────────────────────────────────────

// ─── DOM ──────────────────────────────────────────────────────────────────────
const views = {
  welcome: document.getElementById("view-welcome"),
  profile: document.getElementById("view-profile"),
  loading: document.getElementById("view-loading"),
};
const inputName = document.getElementById("input-name");
const btnStart = document.getElementById("btn-start");
const errorMsg = document.getElementById("error-msg");
const avatarEl = document.getElementById("avatar-initials");
const profileName = document.getElementById("profile-name");
const tokenValue = document.getElementById("token-value");
const btnCopy = document.getElementById("btn-copy");
const notifArea = document.getElementById("notif-area");
const loadingMsg = document.getElementById("loading-msg");
const btnLogout = document.getElementById("btn-logout");

// ─── State ────────────────────────────────────────────────────────────────────
let messaging = null;

// ─── View switcher ────────────────────────────────────────────────────────────
function showView(name) {
  Object.values(views).forEach((v) => v.classList.remove("active"));
  views[name].classList.add("active");
}

// ─── Initials helper ──────────────────────────────────────────────────────────
function initials(name) {
  return name
    .trim()
    .split(/\s+/)
    .map((w) => w[0].toUpperCase())
    .slice(0, 2)
    .join("");
}

// ─── Main: user taps Go ───────────────────────────────────────────────────────
btnStart.addEventListener("click", startFlow);
inputName.addEventListener("keydown", (e) => {
  if (e.key === "Enter") startFlow();
});

async function startFlow() {
  const name = inputName.value.trim();
  if (!name) {
    showError("Please enter your name.");
    return;
  }

  hideError();
  showView("loading");
  setLoading("Checking browser support…");

  try {
    // ── 1. Check support ─────────────────────────────────────────────────────
    if (!("Notification" in window))
      throw new Error("This browser doesn't support notifications.");
    if (!("serviceWorker" in navigator))
      throw new Error("Service workers not supported.");

    // ── 2. Init Firebase (only once) ─────────────────────────────────────────
    if (!firebase.apps.length) {
      firebase.initializeApp(FIREBASE_CONFIG);
    }
    messaging = firebase.messaging();
    console.log(messaging);

    // ── 3. Register service worker ───────────────────────────────────────────
    setLoading("Registering service worker…");
    const swReg = await navigator.serviceWorker.register("/service-worker.js");
    console.log(swReg);

    // ── 4. Request permission + get token ────────────────────────────────────
    //    This triggers the browser "Allow notifications?" popup.
    setLoading("Allow notifications when the browser asks…");

    const token = await messaging.getToken({
      vapidKey: VAPID_KEY,
      serviceWorkerRegistration: swReg,
    });

    console.log("token", token);
    if (!token) throw new Error("Permission denied or token unavailable.");

    // ── 5. Optionally save to backend ────────────────────────────────────────
    if (BACKEND_SAVE_URL) {
      setLoading("Saving to server…");
      await fetch(`${BACKEND_SAVE_URL}/fcmToken?token=${token} && id=${name}`, {
        method: "POST",
      }).catch(() => {}); // non-fatal
    }

    // ── 6. Listen for foreground messages ────────────────────────────────────
    messaging.onMessage((payload) => {
      showInAppNotification(
        payload.notification?.title || "New message",
        payload.notification?.body || "",
      );
    });

    // ── 7. Show profile ──────────────────────────────────────────────────────
    avatarEl.textContent = initials(name);
    profileName.textContent = name;
    tokenValue.textContent = token;
    showView("profile");
  } catch (err) {
    showView("welcome");

    // Translate common errors into plain English
    if (err.code === "messaging/permission-blocked") {
      showError(
        "Notifications are blocked. Open browser settings and allow notifications for this site.",
      );
    } else if (err.code === "messaging/unsupported-browser") {
      showError("Your browser doesn't support push notifications.");
    } else if (err.message?.includes("push service error")) {
      showError(
        "Push service error — make sure the VAPID key matches your Firebase project and the page is served over HTTPS (or localhost).",
      );
    } else {
      showError(err.message || "Something went wrong.");
    }
  }
}

// ─── Copy token ───────────────────────────────────────────────────────────────
btnCopy.addEventListener("click", () => {
  const text = tokenValue.textContent;
  if (!text || text === "—") return;
  navigator.clipboard.writeText(text).then(() => {
    btnCopy.querySelector("span").textContent = "Copied!";
    setTimeout(
      () => (btnCopy.querySelector("span").textContent = "Copy"),
      2000,
    );
  });
});

// ─── Start over ───────────────────────────────────────────────────────────────
btnLogout.addEventListener("click", () => {
  inputName.value = "";
  notifArea.innerHTML = "";
  showView("welcome");
});

// ─── In-app notification banner ───────────────────────────────────────────────
function showInAppNotification(title, body) {
  const el = document.createElement("div");
  el.className = "notif-toast";
  el.innerHTML = `<div class="nt-title">🔔 ${title}</div>${body ? `<div class="nt-body">${body}</div>` : ""}`;
  notifArea.prepend(el);
}

// ─── Helpers ──────────────────────────────────────────────────────────────────
function setLoading(msg) {
  loadingMsg.textContent = msg;
}
function showError(msg) {
  errorMsg.textContent = msg;
  errorMsg.classList.remove("hidden");
}
function hideError() {
  errorMsg.classList.add("hidden");
}
