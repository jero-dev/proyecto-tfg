var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

// Get an specific game with their prices updated
app.MapGet("/product/{productName}", () => { });

// Create a new register for a video-game
app.MapPost("/product", () => { });

// Update info of a product (like name or broken link) by the user
app.MapPut("/product/{productId}", () => { });

app.Run();