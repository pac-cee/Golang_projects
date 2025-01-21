import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Box,
  IconButton,
  List,
  ListItem,
  ListItemText,
} from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';
import { Category, CategoryType } from '../../types';

interface CategoryFormProps {
  open: boolean;
  onClose: () => void;
  category?: Category | null;
  onSubmit: (data: Omit<Category, 'id' | 'createdAt' | 'updatedAt'>) => Promise<void>;
}

const initialFormData: Omit<Category, 'id' | 'createdAt' | 'updatedAt'> = {
  name: '',
  type: 'expense',
  subcategories: [],
};

export const CategoryForm: React.FC<CategoryFormProps> = ({
  open,
  onClose,
  category,
  onSubmit,
}) => {
  const [formData, setFormData] = useState<Omit<Category, 'id' | 'createdAt' | 'updatedAt'>>(
    initialFormData
  );
  const [newSubcategory, setNewSubcategory] = useState('');

  useEffect(() => {
    if (category) {
      setFormData({
        name: category.name,
        type: category.type,
        subcategories: [...category.subcategories],
      });
    } else {
      setFormData(initialFormData);
    }
  }, [category]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await onSubmit(formData);
    onClose();
  };

  const handleAddSubcategory = () => {
    if (newSubcategory.trim()) {
      setFormData({
        ...formData,
        subcategories: [...formData.subcategories, newSubcategory.trim()],
      });
      setNewSubcategory('');
    }
  };

  const handleRemoveSubcategory = (index: number) => {
    setFormData({
      ...formData,
      subcategories: formData.subcategories.filter((_, i) => i !== index),
    });
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <form onSubmit={handleSubmit}>
        <DialogTitle>{category ? 'Edit Category' : 'Add Category'}</DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, pt: 2 }}>
            <TextField
              label="Name"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              required
            />

            <FormControl fullWidth>
              <InputLabel>Type</InputLabel>
              <Select
                value={formData.type}
                onChange={(e) => setFormData({ ...formData, type: e.target.value as CategoryType })}
                label="Type"
                required
              >
                <MenuItem value="expense">Expense</MenuItem>
                <MenuItem value="income">Income</MenuItem>
              </Select>
            </FormControl>

            <Box>
              <Box sx={{ display: 'flex', gap: 1, mb: 2 }}>
                <TextField
                  label="Add Subcategory"
                  value={newSubcategory}
                  onChange={(e) => setNewSubcategory(e.target.value)}
                  fullWidth
                />
                <Button
                  variant="contained"
                  onClick={handleAddSubcategory}
                  disabled={!newSubcategory.trim()}
                >
                  Add
                </Button>
              </Box>

              <List>
                {formData.subcategories.map((subcategory, index) => (
                  <ListItem
                    key={index}
                    secondaryAction={
                      <IconButton
                        edge="end"
                        aria-label="delete"
                        onClick={() => handleRemoveSubcategory(index)}
                      >
                        <DeleteIcon />
                      </IconButton>
                    }
                  >
                    <ListItemText primary={subcategory} />
                  </ListItem>
                ))}
              </List>
            </Box>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose}>Cancel</Button>
          <Button type="submit" variant="contained" color="primary">
            {category ? 'Update' : 'Add'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};
