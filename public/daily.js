import { generateQuestionElement, fetchQuestions } from "./shared.js";
let currentQuestionSet = [];
let currentQuestionIndex = 0;
let quizContainer;
let progressBar;
let nextButton;
let progressIndicator;
let mistakes = 0;
let totalQuestions = 0;
const storageKey = "quiz-history";
const questionsAnsweredCorrectly = /* @__PURE__ */ new Set();
let answers = [];
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
  answers.push(answer);
  if (isCorrect) {
    questionsAnsweredCorrectly.add(question.id);
  } else {
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
}

function updateProgressIndicator() {
  if (progressIndicator) {
    progressIndicator.innerHTML = `${questionsAnsweredCorrectly.size}/${totalQuestions} Correct`;
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

  for (let i = 0; i < answers.length; i++) {
    const answer = answers[i];
    const question = initialQuestions.find((q) => q.id === answer.question_id);
    if (i % initialQuestions.length == 0) {
      results += "<br>";
    }
    results +=
      question.question.correct_answer === answer.answer ? "🥜" : "💩️️️️️️";
  }
  // TODO: keep row consistent even by adding answers to old questions

  const date = new Date().toLocaleDateString();
  const shareText = `${date}\n${results.replaceAll("<br>", "\n")}\n\n${
    window.location.href
  }`;

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
    location.reload();
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
    location.reload();
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
