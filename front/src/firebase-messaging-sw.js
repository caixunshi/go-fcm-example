importScripts('https://www.gstatic.com/firebasejs/8.10.0/firebase-app.js');
importScripts('https://www.gstatic.com/firebasejs/8.10.0/firebase-messaging.js');

// Initialize the Firebase app in the service worker by passing in the
// messagingSenderId.
firebase.initializeApp({
    apiKey: "AIzaSyAUVTaipaExXUTGGc7e-A3gUiA3Q8i7O8Y",
    authDomain: "shipment-portal-1f1b3.firebaseapp.com",
    projectId: "shipment-portal-1f1b3",
    storageBucket: "shipment-portal-1f1b3.appspot.com",
    messagingSenderId: "333922912996",
    appId: "1:333922912996:web:18467f5642e6fba00efaf1",
    measurementId: "G-5HNY007WZW"
});

// Retrieve an instance of Firebase Messaging so that it can handle background
// messages.
const messaging = firebase.messaging();

messaging.onBackgroundMessage((payload) => {
    console.log('Received background message ', payload);
    // console.log(self)
    var title = payload.data.title
    var options = {
        body: payload.data.body,
        icon: payload.data.url
    }
    self.registration.showNotification(title,
        options);
});

