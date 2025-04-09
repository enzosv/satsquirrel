import { generateQuestionElement, fetchQuestions } from "./shared.js";
let currentQuestionSet = [];
let currentQuestionIndex = 0;
let initialAnswers = [];
const questionsAnsweredCorrectly = /* @__PURE__ */ new Set();
let quizContainer;
let progressBar;
let nextButton;
let progressIndicator;
let mistakes = 0;
let totalQuestions = 0;
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
  nextButton.textContent = "Again";
}
function handleAnswer(question, option) {
  question.user_answer = option;
  renderCurrentQuestion();
  const isCorrect = option === question.question.correct_answer;
  if (initialAnswers) {
    initialAnswers.push({
      question_id: question.id,
      user_answer: option,
    });
  }
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
    window.location.href = window.location.href;
    return;
  }
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
    quizContainer.innerHTML = `<div class="alert alert-warning">No questions loaded. Cannot start review mode.</div>`;
    progressIndicator.textContent = "Error";
    return;
  }
  updateProgressIndicator();
  currentQuestionIndex = 0;
  renderCurrentQuestion();
});
