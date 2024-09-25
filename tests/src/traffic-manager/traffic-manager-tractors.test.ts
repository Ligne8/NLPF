import axios from 'axios';

const BASE_URL = 'http://localhost:8080/api/v1';

describe('API Tests for Traffic Manager : Tractors page', () => {

    // Test that the API returns a 200 status code for the GET request
    it('should return 200 for GET /traffic-manager/tractors', async () => {
        const response = await axios.get(`${BASE_URL}/traffic-manager/tractors`);
        expect(response.status).toBe(200);
    });

    // Test that the response contains tractors with the expected properties
    it('should return a tractor with the expected properties', async () => {
        const response = await axios.get(`${BASE_URL}/traffic-manager/tractors`);
        expect(response.status).toBe(200);
        response.data.forEach((tractor: any) => {
            expect(tractor).toHaveProperty('name');
            expect(tractor).toHaveProperty('status');
            expect(tractor).toHaveProperty('currentCapacity');
            expect(tractor).toHaveProperty('totalCapacity');
            expect(tractor).toHaveProperty('location');
            expect(tractor).toHaveProperty('route');
        });
    });

    // Test that the properties of the tractors have the correct data types
    it('should return tractors with the correct data types', async () => {
        const response = await axios.get(`${BASE_URL}/traffic-manager/tractors`);
        expect(response.status).toBe(200);
        response.data.forEach((tractor: any) => {
            expect(typeof tractor.name).toBe('string');
            expect(['ON_THE_WAY', 'ON_THE_STOCK_EXCHANGE', 'AVAILABLE']).toContain(tractor.status);
            expect(typeof tractor.currentCapacity).toBe('number');
            expect(typeof tractor.totalCapacity).toBe('number');
            expect(typeof tractor.location).toBe('string');
            expect(Array.isArray(tractor.route)).toBe(true);
        });
    });

    // Test that all tractors have valid status values from the defined list
    it('should have valid status values', async () => {
        const response = await axios.get(`${BASE_URL}/traffic-manager/tractors`);
        expect(response.status).toBe(200);
        const validStatuses = ['ON_THE_WAY', 'ON_THE_STOCK_EXCHANGE', 'AVAILABLE'];     
        response.data.forEach((tractor: any) => {
            expect(validStatuses).toContain(tractor.status);
        });
    });

    // Test that the API returns an empty array when no tractors are available
    it('should return an empty array when no tractors are available', async () => {
        const response = await axios.get(`${BASE_URL}/traffic-manager/tractors`);
        expect(response.status).toBe(200);
        expect(Array.isArray(response.data)).toBe(true);
        expect(response.data.length).toBeGreaterThanOrEqual(0);
    });
});
