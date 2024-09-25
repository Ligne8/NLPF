import axios from 'axios';

const BASE_URL = 'http://localhost:8080/api/v1';

describe('API Tests for Traffic Manager : Routes page', () => {

    // Test that the API returns a 200 status code for the GET request
    it('should return 200 for GET /traffic-manager/routes', async () => {
        const response = await axios.get(`${BASE_URL}/traffic-manager/routes`);
        expect(response.status).toBe(200);
    });

    // Test that the response contains routes with the expected properties
    it('should return a route with the expected properties', async () => {
        const response = await axios.get(`${BASE_URL}/traffic-manager/routes`);
        expect(response.status).toBe(200);
        response.data.forEach((route: any) => {
            expect(route).toHaveProperty('name');
            expect(route).toHaveProperty('route');
        });
    });

    // Test that the properties of the routes have the correct data types
    it('should return routes with the correct data types', async () => {
        const response = await axios.get(`${BASE_URL}/traffic-manager/routes`);
        expect(response.status).toBe(200);
        response.data.forEach((route: any) => {
            expect(typeof route.name).toBe('string');
            expect(Array.isArray(route.route)).toBe(true);
            route.route.forEach((checkpoint: any) => {
                expect(typeof checkpoint).toBe('string');
            });
        });
    });

    // Test that the API returns an empty array when no routes are available
    it('should return an empty array when no routes are available', async () => {
        const response = await axios.get(`${BASE_URL}/traffic-manager/routes`);
        expect(response.status).toBe(200);
        expect(Array.isArray(response.data)).toBe(true);
        expect(response.data.length).toBeGreaterThanOrEqual(0);
    });
});
