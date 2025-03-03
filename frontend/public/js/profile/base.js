document.addEventListener("DOMContentLoaded", function () {
    // Simuler les données utilisateur récupérées depuis une base de données ou une API
    const userData = {
        username: "Axel",
        memberSince: "01/01/2020",
        rank: "Modérateur",
        messages: 256,
        location: "France",
        avatar: "avatar.jpg",
        lastMessages: [
            { title: "Titre du sujet 1", date: "20/02/2025", link: "#" },
            { title: "Titre du sujet 2", date: "18/02/2025", link: "#" },
            { title: "Titre du sujet 3", date: "15/02/2025", link: "#" }
        ]
    };

    // Mettre à jour les éléments HTML avec les données utilisateur
    document.getElementById("username").textContent = userData.username;
    document.getElementById("member-since").textContent = userData.memberSince;
    document.getElementById("rank").textContent = userData.rank;
    document.getElementById("messages").textContent = userData.messages;
    document.getElementById("location").textContent = userData.location;
    document.getElementById("avatar").src = userData.avatar;

    // Afficher les derniers messages
    const messagesList = document.getElementById("messages-list");
    userData.lastMessages.forEach(message => {
        const li = document.createElement("li");
        const a = document.createElement("a");
        a.href = message.link;
        a.textContent = message.title;
        li.appendChild(a);
        li.innerHTML += ` - Posté le ${message.date}`;
        messagesList.appendChild(li);
    });
});
