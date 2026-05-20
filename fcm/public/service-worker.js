// service-worker.js
// Must be at the ROOT of your site (same level as index.html).
// It handles background push messages (when the tab is closed / not focused).

importScripts(
  "https://www.gstatic.com/firebasejs/9.22.0/firebase-app-compat.js",
);
importScripts(
  "https://www.gstatic.com/firebasejs/9.22.0/firebase-messaging-compat.js",
);

// ── Same config as index.js ───────────────────────────────────────────────────
firebase.initializeApp({
  apiKey: "AIzaSyBWNuaDcXkFi0b37acxgkBlVpEPpcudW0Y",
  authDomain: "fcm101-5e35e.firebaseapp.com",
  projectId: "fcm101-5e35e",
  storageBucket: "fcm101-5e35e.firebasestorage.app",
  messagingSenderId: "106012606884",
  appId: "1:106012606884:web:4bf44c9a1cb69c75ed8131",
});

const messaging = firebase.messaging();

// ── Background message handler ────────────────────────────────────────────────
// This fires when a push arrives and the tab is in the background / closed.
// FCM auto-shows the notification from the payload, so you usually don't
// need to call showNotification() yourself — but you can customise it here.
messaging.onBackgroundMessage((payload) => {
  console.log("[SW] Background message:", payload);

  const { title = "Notification", body = "" } = payload.notification ?? {};

  self.registration.showNotification(title, {
    body,
  });
});
