document.addEventListener("DOMContentLoaded", function () {
    // Gestion de l'image de profil
    const profileImg = document.getElementById("profile-img");
    const profileImgInput = document.getElementById("profile-img-input");
    const resetPhotoBtn = document.getElementById("reset-photo");

    profileImgInput.addEventListener("change", function (event) {
        const file = event.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = function (e) {
                profileImg.src = e.target.result;
            };
            reader.readAsDataURL(file);
        }
    });

    resetPhotoBtn.addEventListener("click", function () {
        profileImg.src = "https://bootdey.com/img/Content/avatar/avatar1.png";
        profileImgInput.value = "";
    });

    // Gestion des informations du profil
    const usernameInput = document.getElementById("username");
    const nameInput = document.getElementById("name");
    const bioInput = document.getElementById("bio");
    const birthdayInput = document.getElementById("birthday");
    const countryInput = document.getElementById("country");
    const phoneInput = document.getElementById("phone");
    const websiteInput = document.getElementById("website");
    const saveChangesBtn = document.getElementById("save-changes");

    saveChangesBtn.addEventListener("click", function () {
        const userInfo = {
            username: usernameInput.value.trim(),
            name: nameInput.value.trim(),
            bio: bioInput.value.trim(),
            birthday: birthdayInput.value.trim(),
            country: countryInput.value,
            phone: phoneInput.value.trim(),
            website: websiteInput.value.trim()
        };

        if (!userInfo.username || !userInfo.name) {
            alert("Veuillez remplir tous les champs obligatoires !");
            return;
        }

        alert("Modifications enregistrées !");
    });

    // Gestion du changement de mot de passe
    const newPasswordInput = document.getElementById("new-password");
    const repeatPasswordInput = document.getElementById("repeat-password");
    const updatePasswordBtn = document.getElementById("update-password");

    updatePasswordBtn.addEventListener("click", function () {
        if (newPasswordInput.value !== repeatPasswordInput.value) {
            alert("Les mots de passe ne correspondent pas !");
            return;
        }

        alert("Mot de passe mis à jour avec succès !");
    });

    // Gestion des liens sociaux
    const socialInputs = document.querySelectorAll(".social-link");
    socialInputs.forEach(input => {
        input.addEventListener("change", function () {
            alert("Lien social mis à jour : " + this.value);
        });
    });

    // Gestion des préférences de notifications
    const notificationCheckboxes = document.querySelectorAll(".notification-setting");
    notificationCheckboxes.forEach(checkbox => {
        checkbox.addEventListener("change", function () {
            alert("Préférence de notification mise à jour : " + (this.checked ? "Activée" : "Désactivée"));
        });
    });
});
