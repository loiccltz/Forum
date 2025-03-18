
document.addEventListener('DOMContentLoaded', function() {
    // Données de l'utilisateur (simulées)
    let userData = {
      name: "",
      firstname: "",
      birthdate: "",
      birthplace: "",
      email: "",
      phone: "",
      address: "",
      profession: "",
      company: ""
    };
  
    const form = document.getElementById('profile-form');
    const editBtn = document.getElementById('edit-btn');
    const saveBtn = document.getElementById('save-btn');
    const themeSelect = document.getElementById('theme-select');
  
    // Fonction pour mettre à jour les informations du profil
    function updateProfile(data) {
      for (const key in data) {
        const element = document.getElementById(key);
        const input = document.getElementById(`${key}-input`);
        if (element && input) {
          element.textContent = data[key];
          input.value = data[key];
        }
      }
  
      // Mise à jour de l'image de profil avec les initiales
      const initials = (data.firstname[0] + data.name[0]).toUpperCase();
      document.getElementById('profile-picture').textContent = initials;
    }
  
    // Fonction pour basculer entre le mode affichage et le mode édition
    function toggleEditMode() {
      const infoValues = document.querySelectorAll('.info-value');
      const inputs = document.querySelectorAll('input');
      
      infoValues.forEach(value => value.classList.toggle('hidden'));
      inputs.forEach(input => input.classList.toggle('hidden'));
      
      editBtn.classList.toggle('hidden');
      saveBtn.classList.toggle('hidden');
    }
  
    // Gestionnaire d'événement pour le bouton "Modifier"
    editBtn.addEventListener('click', toggleEditMode);
  
    // Gestionnaire d'événement pour le formulaire
    form.addEventListener('submit', function(e) {
      e.preventDefault();
      
      // Mettre à jour les données utilisateur avec les nouvelles valeurs
      for (const key in userData) {
        const input = document.getElementById(`${key}-input`);
        if (input) {
          userData[key] = input.value;
        }
      }
      
      // Mettre à jour l'affichage
      updateProfile(userData);
      toggleEditMode();
      
      // Ici, vous pourriez envoyer les données mises à jour au serveur
      console.log('Données mises à jour:', userData);
    });
  
    // Fonction pour changer le thème
    function changeTheme(theme) {
      const root = document.documentElement;
      switch (theme) {
        case 'light':
          root.style.setProperty('--primary-color', '#3498db');
          root.style.setProperty('--secondary-color', '#2ecc71');
          root.style.setProperty('--background-color', '#ecf0f1');
          root.style.setProperty('--text-color', '#2c3e50');
          root.style.setProperty('--border-color', '#bdc3c7');
          root.style.setProperty('--container-bg', '#ffffff');
          break;
        default:
          root.style.setProperty('--primary-color', '#3498db');
          root.style.setProperty('--secondary-color', '#2ecc71');
          root.style.setProperty('--background-color', '#1b2838');
          root.style.setProperty('--text-color', '#ffffff');
          root.style.setProperty('--border-color', '#bdc3c7');
          root.style.setProperty('--container-bg', '171a21');
      }
    }
  
    // Gestionnaire d'événement pour le sélecteur de thème
    themeSelect.addEventListener('change', function() {
      changeTheme(this.value);
    });
  
    // Initialiser l'affichage du profil
    updateProfile(userData);
  });
  