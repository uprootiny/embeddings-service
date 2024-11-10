#!/usr/bin/env bb

(require '[babashka.process :refer [sh]]
         '[clojure.java.io :as io]
         '[clojure.string :as str]
         '[hiccup.core :refer [html]]
         '[hiccup.page :refer [html5 include-css]]
         '[cheshire.core :as json])

;; Function to execute shell commands
(defn run-command [cmd]
  (try
    (let [{:keys [out err exit]} (sh "bash" "-c" cmd)]
      (if (zero? exit)
        (str/trim out)
        (str "Error executing command: " err)))
    (catch Exception e
      (str "Exception: " (.getMessage e)))))

;; Collect system information
(defn get-system-info []
  {:hostname (run-command "hostname")
   :os (run-command "lsb_release -d | cut -f2")
   :uptime (run-command "uptime -p")
   :kernel (run-command "uname -r")
   :architecture (run-command "uname -m")})

;; Collect active network services
(defn get-network-services []
  (run-command "sudo lsof -i -P -n | grep LISTEN | awk '{print $1, $9}' | sort | uniq"))

;; Collect recent projects
(defn get-recent-projects []
  (let [project-paths ["/home/$USER/ClojureProjects" "/home/$USER/tinystatus" "/home/$USER/financial-dashboard"]
        cmd (str "find " (str/join " " project-paths)
                 " -maxdepth 1 -type d -printf '%TY-%Tm-%Td %TT %p\n' | sort -r | head -n 5")]
    (map (fn [line]
           (let [[date time path] (str/split line #" " 3)]
             {:name (last (str/split path #"/"))
              :last-modified (str date " " time)}))
         (str/split-lines (run-command cmd)))))

;; Collect scraper statuses
(defn get-scraper-statuses []
  ;; Mock data for scrapers
  [{:name "News Scraper" :status "Running"}
   {:name "Stock Scraper" :status "Idle"}
   {:name "Event Correlation Scraper" :status "Stopped"}])

;; Collect available analyses
(defn get-analyses []
  ["Sentiment Analysis" "Stock Price Correlation" "Market Trend Forecast"])

;; Collect embeddings and intent mappings
(defn get-embeddings-mappings []
  ;; Mock data for embeddings mappings
  [{:intent "Scrape Financial News"
    :project "news_scraper"
    :params "news_params.json"}
   {:intent "Perform Sentiment Analysis"
    :project "sentiment_analyzer"
    :params "sentiment_params.json"}])

;; Generate the HTML content using Hiccup
(defn generate-html [data]
  (html5
   {:lang "en"}
   [:head
    [:meta {:charset "UTF-8"}]
    [:title "Financial Intelligence Dashboard"]
    ;; Include styles.css if available
    (when (.exists (io/file "styles.css"))
      (include-css "styles.css"))
    ;; PWA Manifest and Service Worker Registration
    [:link {:rel "manifest" :href "manifest.json"}]
    [:script
     "if ('serviceWorker' in navigator) {
        navigator.serviceWorker.register('service-worker.js')
          .then(function(registration) {
            console.log('Service Worker registered with scope:', registration.scope);
          }).catch(function(error) {
            console.log('Service Worker registration failed:', error);
          });
      }"]]
   [:body
    [:header
     [:h1 "Financial Intelligence Dashboard"]]
    [:nav
     [:button {:onclick "showSection('overview')"} "Overview"]
     [:button {:onclick "showSection('projects')"} "Projects"]
     [:button {:onclick "showSection('scrapers')"} "Scrapers"]
     [:button {:onclick "showSection('analysis')"} "Analysis"]
     [:button {:onclick "showSection('embeddings')"} "Embeddings"]
     [:button {:onclick "showSection('console')"} "Console"]]
    [:main
     ;; Overview Section
     [:section {:id "overview" :class "section active"}
      [:h2 "System Overview"]
      [:div {:class "card"}
       [:h3 "System Information"]
       [:p [:strong "Hostname: "] (:hostname data)]
       [:p [:strong "OS: "] (:os data)]
       [:p [:strong "Uptime: "] (:uptime data)]
       [:p [:strong "Kernel: "] (:kernel data)]
       [:p [:strong "Architecture: "] (:architecture data)]]
      [:div {:class "card"}
       [:h3 "Active Network Services"]
       (for [line (str/split-lines (:network-services data))]
         [:p line])]]
     ;; Projects Section
     [:section {:id "projects" :class "section"}
      [:h2 "Projects"]
      [:div {:class "card"}
       [:h3 "Recently Modified Projects"]
       [:table
        [:thead
         [:tr
          [:th "Project"]
          [:th "Last Modified"]
          [:th "Actions"]]]
        [:tbody
         (for [proj (:recent-projects data)]
           [:tr
            [:td (:name proj)]
            [:td (:last-modified proj)]
            [:td
             [:button {:class "primary" :onclick (str "viewProject('" (:name proj) "')")} "Details"]]])]]]]
     ;; Scrapers Section
     [:section {:id "scrapers" :class "section"}
      [:h2 "Scraper Management"]
      [:div {:class "card"}
       [:h3 "Active Scrapers"]
       [:ul
        (for [scraper (:scraper-statuses data)]
          [:li
           [:strong (:name scraper)] " - "
           [:span {:class (if (= "Running" (:status scraper)) "highlight" "")} (:status scraper)]
           [:div {:class "actions"}
            [:button {:onclick (str "viewLogs('" (str/lower-case (str/replace (:name scraper) " " "-")) "')")} "View Logs"]
            [:button {:onclick (str "editConfig('" (str/lower-case (str/replace (:name scraper) " " "-")) "')")} "Edit Config"]
            [:button {:onclick (str "controlScraper('" (str/lower-case (str/replace (:name scraper) " " "-")) "', '" (if (= "Running" (:status scraper)) "stop" "start") "')")} (if (= "Running" (:status scraper)) "Stop" "Start")]]])]]]
     ;; Analysis Section
     [:section {:id "analysis" :class "section"}
      [:h2 "Data Analysis"]
      [:div {:class "card"}
       [:h3 "Available Analyses"]
       [:ul
        (for [analysis (:analyses data)]
          [:li
           [:strong analysis]
           [:div {:class "actions"}
            [:button {:onclick (str "runAnalysis('" (str/lower-case (str/replace analysis " " "-")) "')")} "Run Analysis"]]])]]]
     ;; Embeddings Section
     [:section {:id "embeddings" :class "section"}
      [:h2 "Embeddings Service"]
      [:div {:class "card"}
       [:h3 "Intent Mapping"]
       [:p "Enter your intent, and we'll map it to the appropriate projects and tasks."]
       [:input {:type "text" :id "user-intent" :placeholder "e.g., Analyze the impact of recent tech news on stock prices"}]
       [:div {:class "actions"}
        [:button {:class "primary" :onclick "mapIntent()"} "Map Intent"]]
       [:div {:id "intent-results"}
        ;; This would be populated dynamically via JavaScript
        ]]]
     ;; Console Section
     [:section {:id "console" :class "section"}
      [:h2 "Console"]
      [:div {:id "console"}]
      [:input {:type "text" :id "console-input" :placeholder "Enter command..." :onkeydown "if(event.key === 'Enter') executeCommand()"}]]]
    [:footer
     [:p "&copy; 2024 Financial Intelligence Dashboard"]]
    ;; JavaScript functions
    [:script
     "function showSection(sectionId) {
        document.querySelectorAll('.section').forEach(section => {
          section.classList.remove('active');
          if (section.id === sectionId) {
            section.classList.add('active');
          }
        });
      }

      function viewProject(projectName) {
        alert('Viewing details for project: ' + projectName);
      }

      function viewLogs(scraperName) {
        alert('Viewing logs for: ' + scraperName);
      }

      function editConfig(scraperName) {
        alert('Editing config for: ' + scraperName);
      }

      function controlScraper(scraperName, action) {
        alert(action.charAt(0).toUpperCase() + action.slice(1) + 'ing scraper: ' + scraperName);
      }

      function runAnalysis(analysisType) {
        alert('Running ' + analysisType + ' analysis...');
      }

      function mapIntent() {
        const intent = document.getElementById('user-intent').value.trim();
        if (!intent) {
          alert('Please enter an intent.');
          return;
        }
        // Simulate intent mapping
        const intentResults = document.getElementById('intent-results');
        intentResults.innerHTML = '<p><strong>Intent:</strong> ' + intent + '</p>' +
                                  '<p><strong>Mapped Project:</strong> Analysis Module</p>' +
                                  '<p><strong>Suggested Actions:</strong></p>' +
                                  '<ul><li>Run sentiment analysis</li><li>Generate market report</li></ul>';
      }

      function executeCommand() {
        const consoleInput = document.getElementById('console-input');
        const command = consoleInput.value.trim();
        if (!command) return;
        const consoleOutput = document.getElementById('console');
        consoleOutput.innerHTML += '<div>$ ' + command + '</div><div>Command executed successfully.</div>';
        consoleInput.value = '';
        consoleOutput.scrollTop = consoleOutput.scrollHeight;
      }"]]))

;; Main function to collect data and generate HTML
(defn -main []
  (let [system-info (get-system-info)
        network-services (get-network-services)
        recent-projects (get-recent-projects)
        scraper-statuses (get-scraper-statuses)
        analyses (get-analyses)
        embeddings-mappings (get-embeddings-mappings)
        data {:hostname (:hostname system-info)
              :os (:os system-info)
              :uptime (:uptime system-info)
              :kernel (:kernel system-info)
              :architecture (:architecture system-info)
              :network-services network-services
              :recent-projects recent-projects
              :scraper-statuses scraper-statuses
              :analyses analyses
              :embeddings-mappings embeddings-mappings}]
    ;; Generate HTML content
    (spit "dashboard.html" (generate-html data))
    (println "Dashboard generated: dashboard.html")))

;; Run the main function
(-main)
