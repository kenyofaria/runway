import React, { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import './AppReviews.css';

const AppReviews = ({ appId, appName, onBack, selectedHours, setSelectedHours }) => {
    const [reviews, setReviews] = useState([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState(null);
    const [searchParams, setSearchParams] = useSearchParams();

    // Sync URL search params with selectedHours state
    useEffect(() => {
        const hoursFromUrl = searchParams.get('hours');
        if (hoursFromUrl && hoursFromUrl !== selectedHours.toString()) {
            setSelectedHours(hoursFromUrl === 'all' ? 'all' : parseInt(hoursFromUrl));
        }
    }, [searchParams, selectedHours, setSelectedHours]);

    useEffect(() => {
        if (appId) {
            const fetchReviews = async () => {
                setIsLoading(true);
                setError(null);
                try {
                    let url = `${process.env.REACT_APP_API_URL}/app/reviews?id=${appId}`;
                    if (selectedHours !== 'all') {
                        url += `&hours=${selectedHours}`;
                    }

                    const response = await fetch(url);
                    if (!response.ok) {
                        throw new Error('Network response was not ok');
                    }
                    const data = await response.json();
                    setReviews(data);
                } catch (err) {
                    setError(err.message);
                } finally {
                    setIsLoading(false);
                }
            };
            fetchReviews();
        }
    }, [appId, selectedHours]);

    const handleHoursChange = (event) => {
        const newHours = event.target.value;
        setSelectedHours(newHours);

        // Update URL search params
        const newSearchParams = new URLSearchParams(searchParams);
        if (newHours === 'all') {
            newSearchParams.set('hours', 'all');
        } else {
            newSearchParams.set('hours', newHours);
        }
        setSearchParams(newSearchParams);
    };

    return (
        <div className="reviews-container">
            <button onClick={onBack} className="back-button">Back to App List</button>
            <h2>Reviews for {appName}</h2>

            <div className="filter-container">
                <label htmlFor="hours-filter">Show reviews from the last:</label>
                <select id="hours-filter" value={selectedHours} onChange={handleHoursChange}>
                    <option value="24">24 hours</option>
                    <option value="48">48 hours</option>
                    <option value="72">72 hours</option>
                    <option value="96">96 hours</option>
                    <option value="all">All</option>
                </select>
            </div>

            {isLoading && <p>Loading reviews...</p>}
            {error && <p className="error-message">Error: {error}</p>}

            {reviews && reviews.length > 0 ? (
                <div className="reviews-list">
                    {reviews.map((review) => (
                        <div key={review.id} className="review-card">
                            <p><strong>Author:</strong> {review.author}</p>
                            <p><strong>Score:</strong> {review.score} / 5</p>
                            <p>{review.content}</p>
                            <p className="review-time">Time: {new Date(review.time).toLocaleString()}</p>
                        </div>
                    ))}
                </div>
            ) : (
                !isLoading && (
                    <p>
                        No reviews found for this app {selectedHours === 'all' ? '' : `in the last ${selectedHours} hours.`}
                    </p>
                )
            )}
        </div>
    );
};

export default AppReviews;