import {
  generateQuestionElement,
  fetchQuestions,
  storageKey,
} from "./shared.js";

// ui elements
let quizContainer;
let progressBar;
let nextButton;
let progressIndicator;

// quiz state
let currentQuestionSet = [];
let currentQuestionIndex = 0;
const questionsAnsweredCorrectly = new Set();
let attempts = [[]];
let currentAttemptIndex = 0;
let initialQuestions = [];

function renderCurrentQuestion() {
  if (!quizContainer || currentQuestionIndex >= currentQuestionSet.length) {
    console.error(
      "Cannot render question: container missing or index out of bounds."
    );
    return;
  }
  const question = currentQuestionSet[currentQuestionIndex];
  quizContainer.innerHTML = "";
  const div = generateQuestionElement(
    question,
    currentQuestionIndex,
    (option) => {
      handleAnswer(question, option);
    }
  );
  quizContainer.appendChild(div);

  updateNextButtonState(false);
}
function updateNextButtonState(enabled, text) {
  if (!nextButton) return;
  nextButton.disabled = !enabled;
  if (text) {
    nextButton.textContent = text;
    return;
  }
  const isLastQuestionInSet =
    currentQuestionIndex >= currentQuestionSet.length - 1;
  if (!isLastQuestionInSet) {
    nextButton.textContent = "Next";
    return;
  }
  const needsReview = currentQuestionSet.some(
    (q) => !questionsAnsweredCorrectly.has(q.id)
  );
  if (needsReview) {
    nextButton.textContent = "Review";
    return;
  }
  nextButton.textContent = "Done";
}

function handleAnswer(question, option) {
  question.user_answer = option;
  renderCurrentQuestion();
  const isCorrect = option === question.question.correct_answer;
  const answer = {
    question_id: question.id,
    answer: option,
  };
  attempts[currentAttemptIndex].push(answer);
  if (isCorrect) {
    questionsAnsweredCorrectly.add(question.id);
  }
  if (progressBar) {
    progressBar.value = (currentQuestionIndex + 1) / currentQuestionSet.length;
  }
  updateProgressIndicator();
  updateNextButtonState(true);
}

function submitAnswers() {
  const data = localStorage.getItem(storageKey);
  const history = data ? JSON.parse(data) : [];
  const attempt = {
    timestamp: new Date().toISOString(),
    answers: attempts[0],
  };

  history.push(attempt);
  localStorage.setItem(storageKey, JSON.stringify(history));
}

function updateProgressIndicator() {
  if (!progressIndicator) {
    return;
  }
  progressIndicator.innerHTML = emojify(
    attempts[currentAttemptIndex],
    currentQuestionSet
  );
}

function emojify(answers, set) {
  let results = "";
  for (const question of set) {
    const answer = answers.find((a) => a.question_id === question.id);
    if (!answer) {
      results += "🐿️";
      continue;
    }
    results +=
      question.question.correct_answer === answer.answer ? "🥜" : "💩️️️️️️";
  }
  return results;
}

function calculateStreak(entries) {
  // entries: [{ date: 'YYYY-MM-DD', value: number }, ...]
  const datesWithData = new Set(entries.map((entry) => entry.date));

  let streak = 0;
  let today = new Date();

  while (true) {
    const yyyyMmDd = today.toISOString().split("T")[0];
    if (datesWithData.has(yyyyMmDd)) {
      streak++;
      today.setDate(today.getDate() - 1);
    } else {
      break;
    }
  }

  return streak;
}

function formatDataForHeatmap(history) {
  const dateMap = {};
  let earliest = "9999-12-31";

  for (const entry of history) {
    const date = entry.timestamp.split("T")[0]; // 'YYYY-MM-DD'
    dateMap[date] = (dateMap[date] || 0) + 1;
    if (date < earliest) {
      earliest = date;
    }
  }

  const formatted = Object.entries(dateMap).map(([date, value]) => ({
    date,
    value,
  }));
  return {
    data: formatted,
    earliest: new Date(earliest),
  };
}

function showComplete() {
  const data = localStorage.getItem(storageKey);
  const history = data ? JSON.parse(data) : [];
  const heatmapData = formatDataForHeatmap(history);
  const streak = calculateStreak(heatmapData.data);

  let results = "";
  for (let i = 0; i < attempts.length; i++) {
    const answers = attempts[i];
    results += emojify(answers, initialQuestions) + "<br>";
  }

  const date = new Date().toLocaleDateString();
  const shareText = `${date}\n🔥 ${streak} day streak\n${results.replaceAll(
    "<br>",
    "\n"
  )}\n${window.location.href}`;

  // Replace the page contents
  const container = document.getElementById("container");
  if (!container) return;

  container.innerHTML = `
    <div class="complete-screen">
      <h2>🔥 ${streak} day streak 🔥</h2>
      <div class="result-grid">
        ${results}
      </div>
      <button class="share-button">Share</button>
      <p>Check back tomorrow for more questions<br>or checkout <a href=https://www.opensat.fun>OpenSAT</a> to explore further</p>
      <div id="cal-heatmap"></div>
      <div class="copied-toast">Results copied to clipboard!</div>
    </div>
  `;

  const cal = new CalHeatmap();
  cal.paint({
    animationDuration: 0,
    itemSelector: "#cal-heatmap",
    domain: { type: "month" },
    subDomain: { type: "day", radius: 2 },
    data: { source: heatmapData.data, x: "date", y: "value" },
    date: { start: heatmapData.earliest },
  });

  const shareButton = container.querySelector(".share-button");
  const toast = container.querySelector(".copied-toast");

  const showToast = () => {
    toast.style.display = "block";
    setTimeout(() => {
      toast.style.display = "none";
    }, 1500);
  };

  shareButton.addEventListener("click", async () => {
    try {
      await navigator.clipboard.writeText(shareText);
      showToast();
    } catch (err) {
      console.error("Failed to copy:", err);
    }
  });
}

function nextStep() {
  currentQuestionSet[currentQuestionIndex].user_answer = void 0;
  currentQuestionIndex++;
  if (currentQuestionIndex < currentQuestionSet.length) {
    renderCurrentQuestion();
    updateNextButtonState(false);
    return;
  }

  const questionsToReview = currentQuestionSet.filter(
    (q) => !questionsAnsweredCorrectly.has(q.id)
  );
  if (questionsToReview.length < 1) {
    showComplete();
    return;
  }
  submitAnswers();
  currentQuestionSet = questionsToReview;
  currentQuestionIndex = 0;
  currentAttemptIndex++;
  attempts.push([]); // Start a new attempt array for the review
  if (progressBar) {
    progressBar.value = 0;
  }
  renderCurrentQuestion();
  updateNextButtonState(false, "Next");
  updateProgressIndicator();
}

document.addEventListener("DOMContentLoaded", async () => {
  currentQuestionSet = await fetchQuestions();
  initialQuestions = [...currentQuestionSet];
  quizContainer = document.getElementById("quiz-container");
  nextButton = document.getElementById("next-button");
  progressIndicator = document.getElementById("progress-indicator");
  progressBar = document.getElementById("progress");
  if (!quizContainer || !nextButton || !progressIndicator) {
    console.error("Required DOM elements not found!");
    return;
  }
  nextButton.addEventListener("click", nextStep);
  if (currentQuestionSet.length < 1) {
    quizContainer.innerHTML = `<div class="alert alert-warning">No questions loaded.</div>`;
    progressIndicator.textContent = "Error";
    return;
  }
  updateProgressIndicator();
  currentQuestionIndex = 0;
  renderCurrentQuestion();
});
