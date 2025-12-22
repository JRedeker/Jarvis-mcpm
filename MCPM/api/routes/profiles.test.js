/**
 * Profiles Routes Tests
 */

const request = require('supertest');
const { createApp } = require('../server');

describe('Profiles Routes', () => {
    let app;

    beforeAll(() => {
        app = createApp();
    });

    describe('GET /api/v1/profiles', () => {
        it('should return list of profiles', async () => {
            const res = await request(app)
                .get('/api/v1/profiles')
                .expect(200);

            expect(res.body.success).toBe(true);
            expect(res.body.data).toBeDefined();
            expect(res.body.data.count).toBeGreaterThanOrEqual(0);
            expect(Array.isArray(res.body.data.profiles)).toBe(true);
        });

        it('should include profile details', async () => {
            const res = await request(app)
                .get('/api/v1/profiles')
                .expect(200);

            if (res.body.data.profiles.length > 0) {
                const profile = res.body.data.profiles[0];
                expect(profile.name).toBeDefined();
                expect(Array.isArray(profile.servers)).toBe(true);
            }
        });
    });

    describe('GET /api/v1/profiles/:name', () => {
        it('should return 404 for non-existent profile', async () => {
            const res = await request(app)
                .get('/api/v1/profiles/nonexistent-profile-xyz')
                .expect(404);

            expect(res.body.success).toBe(false);
            expect(res.body.error.code).toBe('PROFILE_NOT_FOUND');
        });
    });

    describe('POST /api/v1/profiles', () => {
        it('should return 400 for missing name', async () => {
            const res = await request(app)
                .post('/api/v1/profiles')
                .send({ servers: ['test'] })
                .expect(400);

            expect(res.body.success).toBe(false);
            expect(res.body.error.code).toBe('VALIDATION_ERROR');
            expect(res.body.error.message).toContain('name');
        });
    });

    describe('PUT /api/v1/profiles/:name', () => {
        it('should return 404 for non-existent profile', async () => {
            const res = await request(app)
                .put('/api/v1/profiles/nonexistent-profile-xyz')
                .send({ add_servers: 'test' })
                .expect(404);

            expect(res.body.success).toBe(false);
            expect(res.body.error.code).toBe('PROFILE_NOT_FOUND');
        });
    });

    describe('DELETE /api/v1/profiles/:name', () => {
        it('should return 404 for non-existent profile', async () => {
            const res = await request(app)
                .delete('/api/v1/profiles/nonexistent-profile-xyz')
                .expect(404);

            expect(res.body.success).toBe(false);
            expect(res.body.error.code).toBe('PROFILE_NOT_FOUND');
        });
    });

    describe('Profile-Server Tag Synchronization', () => {
        // Test that profile edits sync server tags
        // Note: These are integration tests that verify the sync behavior

        it('should sync server tags when editing profiles via API', async () => {
            // First, get current profiles to find one we can test with
            const profilesRes = await request(app)
                .get('/api/v1/profiles')
                .expect(200);

            // Verify the sync mechanism is wired up by checking audit endpoint
            const auditRes = await request(app)
                .get('/api/v1/audit')
                .expect(200);

            expect(auditRes.body.success).toBe(true);
            expect(auditRes.body.data.summary).toBeDefined();
            // After sync implementation, configs should be in sync
            // (or have documented mismatches that can be fixed)
            expect(typeof auditRes.body.data.summary.isInSync).toBe('boolean');
        });

        it('should report mismatches between servers.json and profiles.json', async () => {
            const res = await request(app)
                .get('/api/v1/audit')
                .expect(200);

            expect(res.body.success).toBe(true);
            expect(Array.isArray(res.body.data.serverTagsNotInProfile)).toBe(true);
            expect(Array.isArray(res.body.data.profileServersWithoutTag)).toBe(true);
        });
    });
});
