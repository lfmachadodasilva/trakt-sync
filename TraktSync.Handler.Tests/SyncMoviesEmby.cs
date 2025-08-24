using AutoFixture;
using Microsoft.Extensions.Logging;
using Moq;
using TraktSync.Emby;
using TraktSync.Emby.Models;
using TraktSync.Trakt;
using TraktSync.Trakt.Models;

namespace TraktSync.Handler.Tests;

public class SyncMoviesEmby
{
    [Fact]
    public async Task SyncMoviesEmby_ShouldMarkEmbyAsWatched()
    {
        var fixture = new Fixture();
        const string imdbId = "imdb1";
        const string embyId = "1";
        var traktClientMock = new Mock<ITraktClient>();
        traktClientMock.Setup(x => x.GetWatchedMoviesAsync())
            .ReturnsAsync(new List<TraktWatchedMoviesResponse>
            {
                new()
                {
                    Movie = new TraktWatchedItemResponse
                    {
                        Ids = new TraktWatchedIdsResponse { Imdb = imdbId },
                        Title = "Movie 1",
                        Year = 2020
                    }
                }
            })
            .Verifiable();
        var embyClientMock = new Mock<IEmbyClient>();
        embyClientMock
            .Setup(x => x.GetMoviesSync())
            .ReturnsAsync(new EmbyResponse
            {
                TotalRecordCount = 1,
                Items = new List<EmbyItemResponse>
                {
                    new()
                    {
                        Name = "Movie 1",
                        Id = embyId,
                        Data = new EmbyItemDataResponse
                        {
                            Played = false
                        },
                        Ids = new EmbyItemIdsResponse { Imdb = imdbId }
                    }
                }
            });
        var loggerMock = new Mock<ILogger<SyncHandler>>();
        fixture.Register(() => loggerMock.Object);
        fixture.Register(() => traktClientMock.Object);
        fixture.Register(() => embyClientMock.Object);
        var syncHandler = fixture.Create<SyncHandler>();

        // Act
        await syncHandler.SyncAsync();
        
        // Assert
        embyClientMock.Verify(
            x => x.MarkAsWatchedAsync(It.Is<string>(y => y == embyId)),
            Times.Once);
        traktClientMock.Verify(
            x => x.MarkAsWatchedAsync(
                It.Is<TraktMarkAsWatchedRequest>(y => y.Movies.Any(z => z.Ids != null && z.Ids.Imdb == imdbId)),
                It.IsAny<bool>()),
            Times.Never);
    }
    
    [Fact]
    public async Task SyncMoviesEmby_ShouldMarkTraktAsWatched()
    {
        var fixture = new Fixture();
        const string imdbId = "imdb1";
        const string embyId = "1";
        var traktClientMock = new Mock<ITraktClient>();
        traktClientMock.Setup(x => x.GetWatchedMoviesAsync())
            .ReturnsAsync(new List<TraktWatchedMoviesResponse>())
            .Verifiable();
        var embyClientMock = new Mock<IEmbyClient>();
        embyClientMock
            .Setup(x => x.GetMoviesSync())
            .ReturnsAsync(new EmbyResponse
            {
                TotalRecordCount = 1,
                Items = new List<EmbyItemResponse>
                {
                    new()
                    {
                        Name = "Movie 1",
                        Id = embyId,
                        Data = new EmbyItemDataResponse
                        {
                            Played = true
                        },
                        Ids = new EmbyItemIdsResponse { Imdb = imdbId }
                    }
                }
            });
        var loggerMock = new Mock<ILogger<SyncHandler>>();
        fixture.Register(() => loggerMock.Object);
        fixture.Register(() => traktClientMock.Object);
        fixture.Register(() => embyClientMock.Object);
        var syncHandler = fixture.Create<SyncHandler>();

        // Act
        await syncHandler.SyncAsync();
        
        // Assert
        embyClientMock.Verify(
            x => x.MarkAsWatchedAsync(It.Is<string>(y => y == embyId)),
            Times.Never);
        traktClientMock.Verify(
            x => x.MarkAsWatchedAsync(
                It.Is<TraktMarkAsWatchedRequest>(y => y.Movies.Any(z => z.Ids != null && z.Ids.Imdb == imdbId)),
                It.IsAny<bool>()),
            Times.Once);
    }
}