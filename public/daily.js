import { generateQuestionElement, fetchQuestions } from "./shared.js";
let currentQuestionSet = [];
let currentQuestionIndex = 0;
let answers = [];
let quizContainer;
let progressBar;
let nextButton;
let progressIndicator;
let mistakes = 0;
let totalQuestions = 0;
const storageKey = "quiz-history";

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
  nextButton.textContent = "Done";
  return;
  // REVIEW MODE
  const needsReview = currentQuestionSet.some(
    (q) => !questionsAnsweredCorrectly.has(q.id)
  );
  if (needsReview) {
    nextButton.textContent = "Done";
    return;
  }
  nextButton.textContent = "Again";
}
function handleAnswer(question, option) {
  question.user_answer = option;
  renderCurrentQuestion();
  answers.push({
    question_id: question.id,
    answer: option,
  });
  if (option != question.question.correct_answer) {
    mistakes++;
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
    answers: answers,
  };

  history.push(attempt);
  localStorage.setItem(storageKey, JSON.stringify(history));

  // globalThis.location.href = `attempt.html?index=${history.length - 1}`;
}

function updateProgressIndicator() {
  if (progressIndicator) {
    progressIndicator.innerHTML = `${
      answers.length - mistakes
    }/${totalQuestions} Correct`;
  }
  if (mistakes > 0) {
    const counter = document.getElementById("mistake-counter");
    if (counter) {
      counter.innerHTML = `${mistakes} Mistake${mistakes == 1 ? "" : "s"}`;
    }
  }
}

function showComplete() {
  let results = "";
  for (const answer of answers) {
    const question = currentQuestionSet.find(
      (q) => q.id === answer.question_id
    );
    results +=
      question.question.correct_answer === answer.answer ? "🥜" : "💩️️️️️️";
  }

  const date = new Date().toLocaleDateString();
  const shareText = `${date}\n${results}\n\n${window.location.href}`;

  const modal = document.createElement("div");
  modal.className = "modal";
  modal.innerHTML = `
    <div class="modal-content">
      <span class="close">&times;</span>
      <h2>Practice Complete!</h2>
      <div class="result-grid">${results}</div>
      <div>
        <button class="share-button">Share</button>
      </div>
      <div class="copied-toast">Results copied to clipboard!</div>
    </div>
  `;

  document.body.appendChild(modal);

  const closeButton = modal.querySelector(".close");
  closeButton.onclick = function () {
    modal.remove();
  };

  const shareButton = modal.querySelector(".share-button");
  const toast = modal.querySelector(".copied-toast");

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

  submitAnswers();
  showComplete();
}

document.addEventListener("DOMContentLoaded", async () => {
  currentQuestionSet = await fetchQuestions();
  quizContainer = document.getElementById("quiz-container");
  nextButton = document.getElementById("next-button");
  progressIndicator = document.getElementById("progress-indicator");
  progressBar = document.getElementById("progress");
  if (!quizContainer || !nextButton || !progressIndicator) {
    console.error("Required DOM elements not found!");
    return;
  }
  nextButton.addEventListener("click", nextStep);
  totalQuestions = currentQuestionSet.length;
  if (currentQuestionSet.length < 1) {
    quizContainer.innerHTML = `<div class="alert alert-warning">No questions loaded.</div>`;
    progressIndicator.textContent = "Error";
    return;
  }
  updateProgressIndicator();
  currentQuestionIndex = 0;
  renderCurrentQuestion();
});
