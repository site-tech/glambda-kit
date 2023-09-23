const main = async function(event) {
  try {
    const response = {
      message: "Hello World!",
    };
    return {
      statusCode: 200,
      body: JSON.stringify(response),
    };
  } catch (error) {
    return {
      statusCode: 500,
      body: JSON.stringify({ error: error.message }),
    };
  }
};

module.exports = { main };
