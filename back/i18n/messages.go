package i18n

// messages: catalogue of the messages returned by the API, by key then by language.
var messages = map[string]map[Lang]string{
	// Common
	"common.invalid_id":     {EN: "Invalid identifier.", FR: "Identifiant invalide."},
	"common.invalid_body":   {EN: "Invalid request body.", FR: "Corps de requête invalide."},
	"common.missing_params": {EN: "Missing parameters.", FR: "Paramètres manquants."},
	"common.internal":       {EN: "An internal error occurred.", FR: "Une erreur interne est survenue."},

	// Login
	"auth.missing_fields":       {EN: "Missing required fields.", FR: "Champs requis manquants."},
	"auth.credentials_required": {EN: "Username and password are required.", FR: "Identifiant et mot de passe requis."},
	"auth.invalid_credentials":  {EN: "Incorrect username or password.", FR: "Login ou mot de passe incorrect."},
	"auth.success":              {EN: "OK, the user is logged in.", FR: "OK, l'utilisateur est connecté."},

	// Items (books and journals)
	"items.search_field_required": {EN: "Please fill in at least one field.", FR: "Veuillez renseigner au moins un champ."},
	"items.search_failed":         {EN: "Error while searching.", FR: "Erreur lors de la recherche."},
	"items.title_required":        {EN: "Title is required.", FR: "Le titre est requis."},
	"items.location_required":     {EN: "Location is required.", FR: "L'emplacement est requis."},
	"items.deposit_negative":      {EN: "Deposit cannot be negative.", FR: "La caution ne peut pas être négative."},
	"authors.load_failed":         {EN: "Unable to load authors.", FR: "Impossible de charger les auteurs."},
	"locations.load_failed":       {EN: "Unable to load locations.", FR: "Impossible de charger les emplacements."},
	"books.title_isbn_required":   {EN: "Title and ISBN are required.", FR: "Le titre et l'ISBN sont requis."},
	"books.author_required":       {EN: "Author required: select an existing author or enter a first and last name.", FR: "Auteur requis : sélectionnez un auteur existant ou saisissez nom et prénom."},
	"books.author_create_failed":  {EN: "Unable to create the author.", FR: "Impossible de créer l'auteur."},
	"books.create_failed":         {EN: "Unable to create the book.", FR: "Impossible de créer le livre."},
	"books.created":               {EN: "Book created successfully.", FR: "Livre créé avec succès."},
	"journals.number_positive":    {EN: "Issue number must be a positive integer.", FR: "Le numéro doit être un entier positif."},
	"journals.date_required":      {EN: "Publication date is required.", FR: "La date de parution est requise."},
	"journals.create_failed":      {EN: "Unable to create the journal.", FR: "Impossible de créer la revue."},
	"journals.created":            {EN: "Journal created successfully.", FR: "Revue créée avec succès."},

	// Copies
	"copies.load_failed":       {EN: "Unable to load copies.", FR: "Impossible de charger les exemplaires."},
	"copies.barcode_required":  {EN: "Barcode is required.", FR: "Code-barre requis."},
	"copies.barcode_exists":    {EN: "This barcode already exists.", FR: "Ce code-barre existe déjà."},
	"copies.barcode_not_found": {EN: "Barcode not found.", FR: "Code-barre introuvable."},
	"copies.already_borrowed":  {EN: "This copy is already on loan.", FR: "Cet exemplaire est déjà emprunté."},
	"copies.create_failed":     {EN: "Unable to create the copy.", FR: "Impossible de créer l'exemplaire."},

	// Loans
	"loans.invalid_user_id":     {EN: "Invalid user_id parameter.", FR: "Paramètre user_id invalide."},
	"loans.load_failed":         {EN: "Unable to retrieve loans.", FR: "Impossible de récupérer les emprunts."},
	"loans.overdue_load_failed": {EN: "Unable to retrieve overdue loans.", FR: "Impossible de récupérer les emprunts en retard."},
	"loans.overdue":             {EN: "You have %d overdue loan(s). Please return them before borrowing.", FR: "Vous avez %d emprunt(s) en retard, veuillez les retourner avant d'emprunter."},
	"loans.already_borrowing":   {EN: "You are already borrowing a copy of this item.", FR: "Vous empruntez déjà un exemplaire de cet ouvrage."},
	"loans.insufficient_funds":  {EN: "Insufficient deposit balance (€%.2f available, €%.2f required).", FR: "Solde insuffisant (%.2f € disponible, %.2f € requis)."},
	"loans.verify_failed":       {EN: "Unable to verify this loan.", FR: "Impossible de vérifier cet emprunt."},
	"loans.create_failed":       {EN: "Unable to register the loan.", FR: "Impossible d'enregistrer l'emprunt."},
	"loans.created":             {EN: "Loan registered successfully.", FR: "Emprunt enregistré avec succès."},
	"loans.not_found":           {EN: "Loan not found.", FR: "Emprunt introuvable."},
	"loans.return_failed":       {EN: "Unable to return this copy.", FR: "Impossible de rendre cet exemplaire."},
	"loans.returned":            {EN: "Copy returned successfully.", FR: "Exemplaire rendu avec succès."},

	// Deposit
	"deposit.invalid_value":  {EN: "Invalid deposit amount.", FR: "Valeur de caution invalide."},
	"deposit.below_borrowed": {EN: "Total deposit cannot be lower than the amount currently borrowed (€%.2f).", FR: "La caution totale ne peut pas être inférieure au montant déjà emprunté (%.2f €)."},
	"deposit.update_failed":  {EN: "Unable to update the deposit.", FR: "Impossible de mettre à jour la caution."},
	"deposit.updated":        {EN: "Deposit updated.", FR: "Caution mise à jour."},

	// Users
	"users.not_found":                {EN: "User not found.", FR: "Utilisateur introuvable."},
	"users.search_criteria_required": {EN: "At least one search criterion is required.", FR: "Au moins un critère de recherche est requis."},
	"users.departments_load_failed":  {EN: "Unable to load departments.", FR: "Impossible de charger les départements."},
	"users.missing_fields":           {EN: "Required fields are missing or invalid.", FR: "Champs requis manquants ou invalides."},
	"users.invalid_phone":            {EN: "Phone number must have 10 digits and start with 0 (e.g. 0601020304).", FR: "Le numéro de téléphone doit contenir 10 chiffres et commencer par 0 (ex : 0601020304)."},
	"users.invalid_date":             {EN: "Invalid date format (expected: YYYY-MM-DD).", FR: "Format de date invalide (attendu : AAAA-MM-JJ)."},
	"users.year_required":            {EN: "Year of study is required for a student.", FR: "L'année d'étude est requise pour un étudiant."},
	"users.department_required":      {EN: "Department is required for a teacher.", FR: "Le département est requis pour un enseignant."},
	"users.account_check_failed":     {EN: "Error while checking the account.", FR: "Erreur lors de la vérification du compte."},
	"users.login_exists":             {EN: "An account already exists with the login “%s”.", FR: "Un compte existe déjà avec le login « %s »."},
	"users.address_failed":           {EN: "Error while saving the address.", FR: "Erreur lors de l'enregistrement de l'adresse."},
	"users.student_create_failed":    {EN: "Error while creating the student.", FR: "Erreur lors de la création de l'étudiant."},
	"users.teacher_create_failed":    {EN: "Error while creating the teacher.", FR: "Erreur lors de la création de l'enseignant."},
	"users.individual_create_failed": {EN: "Error while creating the member.", FR: "Erreur lors de la création du particulier."},
	"users.created":                  {EN: "User created successfully.", FR: "Utilisateur créé avec succès."},
}
