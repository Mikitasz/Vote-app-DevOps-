document.addEventListener("DOMContentLoaded", () => {
  let dogsVotes = 0;
  let catsVotes = 0;

  const dogsVotesElement = document.getElementById("dogs-votes");
  const catsVotesElement = document.getElementById("cats-votes");

  const updateBackground = () => {
    const totalVotes = dogsVotes + catsVotes;
    const dogsPercentage =
      totalVotes === 0 ? 50 : (dogsVotes / totalVotes) * 100;
    const catsPercentage = 100 - dogsPercentage;

    document.body.style.background = `linear-gradient(to right, #ff9a9e ${dogsPercentage}%, #bae1ff ${dogsPercentage}%)`;
  };

  const fetchVotes = () => {
    fetch("/votes")
      .then((response) => response.json())
      .then((data) => {
        dogsVotes = data.dogs || 0;
        catsVotes = data.cats || 0;
        dogsVotesElement.textContent = dogsVotes;
        catsVotesElement.textContent = catsVotes;
        updateBackground();
      })
      .catch((error) => console.error("Error fetching votes:", error));
  };

  document.querySelector(".vote-button.dogs").addEventListener("click", () => {
    fetch("/vote", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ category: "dogs" }),
    }).then((response) => {
      if (response.ok) {
        dogsVotes++;
        dogsVotesElement.textContent = dogsVotes;
        updateBackground();
      }
    });
  });

  document.querySelector(".vote-button.cats").addEventListener("click", () => {
    fetch("/vote", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ category: "cats" }),
    }).then((response) => {
      if (response.ok) {
        catsVotes++;
        catsVotesElement.textContent = catsVotes;
        updateBackground();
      }
    });
  });

  // Fetch initial votes on page load
  fetchVotes();
});
